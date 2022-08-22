package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

type SyncCommitteeRunner struct {
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

func NewSyncCommitteeRunner(
	beaconNetwork types.BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.KeyManager,
	valCheck qbft.ProposedValueCheckF,
) Runner {
	return &SyncCommitteeRunner{
		BeaconRoleType: types.BNRoleSyncCommittee,
		BeaconNetwork:  beaconNetwork,
		Share:          share,
		QBFTController: qbftController,

		beacon:   beacon,
		network:  network,
		signer:   signer,
		valCheck: valCheck,
	}
}

func (r *SyncCommitteeRunner) StartNewDuty(duty *types.Duty) error {
	if err := canStartNewDuty(r, duty); err != nil {
		return err
	}
	r.State = NewRunnerState(r.GetShare().Quorum, duty)
	return r.executeDuty(duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *SyncCommitteeRunner) HasRunningDuty() bool {
	if r.GetState() == nil {
		return false
	}
	return r.GetState().Finished != true
}

func (r *SyncCommitteeRunner) ProcessPreConsensus(signedMsg *SignedPartialSignatureMessage) error {
	return errors.New("no pre consensus sigs required for sync committee role")
}

func (r *SyncCommitteeRunner) ProcessConsensus(signedMsg *qbft.SignedMessage) error {
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
	msg, err := signBeaconObject(r, types.SSZBytes(decidedValue.SyncCommitteeBlockRoot[:]), decidedValue.Duty.Slot, types.DomainSyncCommittee)
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

func (r *SyncCommitteeRunner) ProcessPostConsensus(signedMsg *SignedPartialSignatureMessage) error {
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

		msg := &altair.SyncCommitteeMessage{
			Slot:            r.GetState().DecidedValue.Duty.Slot,
			BeaconBlockRoot: r.GetState().DecidedValue.SyncCommitteeBlockRoot,
			ValidatorIndex:  r.GetState().DecidedValue.Duty.ValidatorIndex,
			Signature:       specSig,
		}
		if err := r.GetBeaconNode().SubmitSyncMessage(msg); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed sync committee")
		}
	}
	r.GetState().Finished = true
	return nil
}

func (r *SyncCommitteeRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return []ssz.HashRoot{}, types.DomainError, errors.New("no expected pre consensus roots for sync committee")
}

// executeDuty steps:
// 1) get sync block root from BN
// 2) start consensus on duty + block root data
// 3) Once consensus decides, sign partial block root and broadcast
// 4) collect 2f+1 partial sigs, reconstruct and broadcast valid sync committee sig to the BN
func (r *SyncCommitteeRunner) executeDuty(duty *types.Duty) error {
	// TODO - waitOneThirdOrValidBlock

	root, err := r.GetBeaconNode().GetSyncMessageBlockRoot()
	if err != nil {
		return errors.Wrap(err, "failed to get sync committee block root")
	}

	input := &types.ConsensusData{
		Duty:                   duty,
		SyncCommitteeBlockRoot: root,
	}

	if err := decide(r, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (r *SyncCommitteeRunner) GetNetwork() Network {
	return r.network
}

func (r *SyncCommitteeRunner) GetBeaconNetwork() types.BeaconNetwork {
	return r.BeaconNetwork
}

func (r *SyncCommitteeRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *SyncCommitteeRunner) GetBeaconRole() types.BeaconRole {
	return r.BeaconRoleType
}

func (r *SyncCommitteeRunner) GetShare() *types.Share {
	return r.Share
}

func (r *SyncCommitteeRunner) GetState() *State {
	return r.State
}

func (r *SyncCommitteeRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *SyncCommitteeRunner) GetQBFTController() *qbft.Controller {
	return r.QBFTController
}

func (r *SyncCommitteeRunner) GetSigner() types.KeyManager {
	return r.signer
}

// Encode returns the encoded struct in bytes or error
func (r *SyncCommitteeRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *SyncCommitteeRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *SyncCommitteeRunner) GetRoot() ([]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}
