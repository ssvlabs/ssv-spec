package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type AggregatorRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewAggregatorRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) Runner {
	return &AggregatorRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType:     types.RoleAggregator,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		valCheck:       valCheck,
	}
}

func (r *AggregatorRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	return r.BaseRunner.baseStartNewDuty(r, duty, quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *AggregatorRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *AggregatorRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing selection proof message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// only 1 root, verified by basePreConsensusMsgProcessing
	root := roots[0]
	// reconstruct selection proof sig
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
	if err != nil {
		// If the reconstructed signature verification failed, fall back to verifying each partial signature
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee,
			r.GetShare().ValidatorIndex)
		return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
	}

	duty := r.GetState().StartingDuty.(*types.ValidatorDuty)

	// TODO waitToSlotTwoThirds

	// get block data
	res, ver, err := r.GetBeaconNode().SubmitAggregateSelectionProof(duty.Slot, duty.CommitteeIndex, duty.CommitteeLength, duty.ValidatorIndex, fullSig)
	if err != nil {
		return errors.Wrap(err, "failed to submit aggregate and proof")
	}

	byts, err := res.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal aggregate and proof")
	}
	input := &types.ValidatorConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, input.Duty.Slot, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (r *AggregatorRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	aggregateAndProof, err := decidedValue.(*types.ValidatorConsensusData).GetAggregateAndProof()
	if err != nil {
		return errors.Wrap(err, "could not get aggregate and proof")
	}

	// specific duty sig
	msg, err := r.BaseRunner.signBeaconObject(r, r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty),
		aggregateAndProof,
		decidedValue.(*types.ValidatorConsensusData).Duty.Slot,
		types.DomainAggregateAndProof)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     decidedValue.(*types.ValidatorConsensusData).Duty.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := postConsensusMsg.Encode()
	if err != nil {
		return err
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	sig, err := r.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{r.operatorSigner.GetOperatorID()},
		SSVMessage:  ssvMsg,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast.SSVMessage.GetID(), msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (r *AggregatorRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
		if err != nil {
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			for _, root := range roots {
				r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PostConsensusContainer, root,
					r.GetShare().Committee, r.GetShare().ValidatorIndex)
			}
			return errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		cd := &types.ValidatorConsensusData{}
		err = cd.Decode(r.GetState().DecidedValue)
		if err != nil {
			return errors.Wrap(err, "could not create consensus data")
		}
		aggregateAndProof, err := cd.GetAggregateAndProof()
		if err != nil {
			return errors.Wrap(err, "could not get aggregate and proof")
		}

		msg := &phase0.SignedAggregateAndProof{
			Message:   aggregateAndProof,
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
	return []ssz.HashRoot{types.SSZUint64(r.GetState().StartingDuty.DutySlot())}, types.DomainSelectionProof, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *AggregatorRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	cd := &types.ValidatorConsensusData{}
	err := cd.Decode(r.GetState().DecidedValue)
	if err != nil {
		return nil, types.DomainError, errors.Wrap(err, "could not create consensus data")
	}
	aggregateAndProof, err := cd.GetAggregateAndProof()
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not get aggregate and proof")
	}

	return []ssz.HashRoot{aggregateAndProof}, types.DomainAggregateAndProof, nil
}

// executeDuty steps:
// 1) sign a partial selection proof and wait for 2f+1 partial sigs from peers
// 2) reconstruct selection proof and send SubmitAggregateSelectionProof to BN
// 3) start consensus on duty + aggregation data
// 4) Once consensus decides, sign partial aggregation data and broadcast
// 5) collect 2f+1 partial sigs, reconstruct and broadcast valid SignedAggregateSubmitRequest sig to the BN
func (r *AggregatorRunner) executeDuty(duty types.Duty) error {
	// sign selection proof
	msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.ValidatorDuty), types.SSZUint64(duty.DutySlot()), duty.DutySlot(),
		types.DomainSelectionProof)
	if err != nil {
		return errors.Wrap(err, "could not sign randao")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.SelectionProofPartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := msgs.Encode()
	if err != nil {
		return err
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	sig, err := r.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{r.operatorSigner.GetOperatorID()},
		SSVMessage:  ssvMsg,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast.SSVMessage.GetID(), msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial selection proof sig")
	}
	return nil
}

func (r *AggregatorRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *AggregatorRunner) GetNetwork() Network {
	return r.network
}

func (r *AggregatorRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *AggregatorRunner) GetShare() *types.Share {
	// there is only 1 share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *AggregatorRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *AggregatorRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *AggregatorRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *AggregatorRunner) GetOperatorSigner() types.OperatorSigner {
	return r.operatorSigner
}
