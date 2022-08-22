package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

type AggregatorRunner struct {
	State          *State
	Share          *types.Share
	QBFTController *qbft.Controller
	BeaconNetwork  types.BeaconNetwork
	BeaconRoleType types.BeaconRole

	beacon   BeaconNode
	network  Network
	signer   types.KeyManager
	valCheck qbft.ProposedValueCheckF
}

func NewAggregatorRunner(
	beaconNetwork types.BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.KeyManager,
	valCheck qbft.ProposedValueCheckF,
) Runner {
	return &AggregatorRunner{
		BeaconRoleType: types.BNRoleAggregator,
		BeaconNetwork:  beaconNetwork,
		Share:          share,
		QBFTController: qbftController,

		beacon:   beacon,
		network:  network,
		signer:   signer,
		valCheck: valCheck,
	}
}

func (r *AggregatorRunner) StartNewDuty(duty *types.Duty) error {
	if err := canStartNewDuty(r, duty); err != nil {
		return err
	}
	r.State = NewRunnerState(r.GetShare().Quorum, duty)
	return r.executeDuty(duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *AggregatorRunner) HasRunningDuty() bool {
	if r.GetState() == nil {
		return false
	}
	return r.GetState().Finished != true
}

func (r *AggregatorRunner) ProcessPreConsensus(signedMsg *SignedPartialSignatureMessage) error {
	quorum, roots, err := basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing randao message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	if len(roots) != 1 {
		return errors.New("too many randao roots")
	}

	root := roots[0]
	// reconstruct selection proof sig
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey)
	if err != nil {
		return errors.Wrap(err, "could not reconstruct selection proof sig")
	}

	duty := r.GetState().StartingDuty

	// TODO waitToSlotTwoThirds

	// get block data
	res, err := r.GetBeaconNode().SubmitAggregateSelectionProof(duty.Slot, duty.CommitteeIndex, fullSig)
	if err != nil {
		return errors.Wrap(err, "failed to submit aggregate and proof")
	}

	input := &types.ConsensusData{
		Duty:              duty,
		AggregateAndProof: res,
	}

	if err := decide(r, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (r *AggregatorRunner) ProcessConsensus(signedMsg *qbft.SignedMessage) error {
	decided, decidedValue, err := baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}
	r.GetState().DecidedValue = decidedValue

	// specific duty sig
	msg, err := signBeaconObject(r, decidedValue.AggregateAndProof, decidedValue.Duty.Slot, types.DomainAggregateAndProof)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}
	postConsensusMsg := &PartialSignatureMessages{
		Type:     PostConsensusPartialSig,
		Messages: []*PartialSignatureMessage{msg},
	}

	postSignedMsg, err := signPostConsensusMsg(r, postConsensusMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign post consensus msg")
	}

	data, err := postSignedMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(r.GetShare().ValidatorPubKey, r.GetBeaconRole()),
		Data:    data,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (r *AggregatorRunner) ProcessPostConsensus(signedMsg *SignedPartialSignatureMessage) error {
	quorum, roots, err := basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus signature")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		msg := &phase0.SignedAggregateAndProof{
			Message:   r.GetState().DecidedValue.AggregateAndProof,
			Signature: specSig,
		}
		if err := r.GetBeaconNode().SubmitSignedAggregateSelectionProof(msg); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed aggregate")
		}
	}
	r.GetState().Finished = true
	return nil
}

func (r *AggregatorRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return []ssz.HashRoot{types.SSZUint64(r.GetState().StartingDuty.Slot)}, types.DomainSelectionProof, nil
}

// executeDuty steps:
// 1) sign a partial selection proof and wait for 2f+1 partial sigs from peers
// 2) reconstruct selection proof and send SubmitAggregateSelectionProof to BN
// 3) start consensus on duty + aggregation data
// 4) Once consensus decides, sign partial aggregation data and broadcast
// 5) collect 2f+1 partial sigs, reconstruct and broadcast valid SignedAggregateSubmitRequest sig to the BN
func (r *AggregatorRunner) executeDuty(duty *types.Duty) error {
	// sign selection proof
	msg, err := signBeaconObject(r, types.SSZUint64(duty.Slot), duty.Slot, types.DomainSelectionProof)
	if err != nil {
		return errors.Wrap(err, "could not sign randao")
	}
	msgs := PartialSignatureMessages{
		Type:     SelectionProofPartialSig,
		Messages: []*PartialSignatureMessage{msg},
	}

	// sign msg
	signature, err := r.GetSigner().SignRoot(msg, types.PartialSignatureType, r.GetShare().SharePubKey)
	if err != nil {
		return errors.Wrap(err, "could not sign PartialSignatureMessage for selection proof")
	}
	signedPartialMsg := &SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: signature,
		Signer:    r.GetShare().OperatorID,
	}

	// broadcast
	data, err := signedPartialMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode selection proof pre-consensus signature msg")
	}
	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(r.GetShare().ValidatorPubKey, r.GetBeaconRole()),
		Data:    data,
	}
	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial selection proof sig")
	}
	return nil
}

func (r *AggregatorRunner) GetNetwork() Network {
	return r.network
}

func (r *AggregatorRunner) GetBeaconNetwork() types.BeaconNetwork {
	return r.BeaconNetwork
}

func (r *AggregatorRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *AggregatorRunner) GetBeaconRole() types.BeaconRole {
	return r.BeaconRoleType
}

func (r *AggregatorRunner) GetShare() *types.Share {
	return r.Share
}

func (r *AggregatorRunner) GetState() *State {
	return r.State
}

func (r *AggregatorRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *AggregatorRunner) GetQBFTController() *qbft.Controller {
	return r.QBFTController
}

func (r *AggregatorRunner) GetSigner() types.KeyManager {
	return r.signer
}

// Encode returns the encoded struct in bytes or error
func (r *AggregatorRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *AggregatorRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *AggregatorRunner) GetRoot() ([]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}
