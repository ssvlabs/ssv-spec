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
	BaseRunner *BaseRunner

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
	highestDecidedSlot phase0.Slot,
) Runner {
	runner := &SyncCommitteeRunner{
		BaseRunner: &BaseRunner{
			BeaconRoleType:     types.BNRoleSyncCommittee,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:   beacon,
		network:  network,
		signer:   signer,
		valCheck: valCheck,
	}

	qbftController.WithCommitExtraLoadManagerF(NewCommitExtraLoadManagerF(runner.BaseRunner, types.BNRoleSyncCommittee, runner.beacon, runner.signer, types.DomainSyncCommittee))

	return runner
}

func (r *SyncCommitteeRunner) StartNewDuty(duty *types.Duty) error {
	return r.BaseRunner.baseStartNewDuty(r, duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *SyncCommitteeRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *SyncCommitteeRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	return errors.New("no pre consensus sigs required for sync committee role")
}

func (r *SyncCommitteeRunner) ProcessConsensus(signedMsg *qbft.SignedMessage) error {
	decided, decidedValue, commitExtraLoadManagerI, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	commitExtraLoadManager := commitExtraLoadManagerI.(*CommitExtraLoadManager)

	roots := commitExtraLoadManager.SigningRoot

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(commitExtraLoadManager.PartialSigContainer, root, r.GetShare().ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus signature")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		blockRoot, err := decidedValue.GetSyncCommitteeBlockRoot()
		if err != nil {
			return errors.Wrap(err, "could not get sync committee block root")
		}

		msg := &altair.SyncCommitteeMessage{
			Slot:            decidedValue.Duty.Slot,
			BeaconBlockRoot: blockRoot,
			ValidatorIndex:  decidedValue.Duty.ValidatorIndex,
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

	root, ver, err := r.GetBeaconNode().GetSyncMessageBlockRoot(duty.Slot)
	if err != nil {
		return errors.Wrap(err, "failed to get sync committee block root")
	}

	input := &types.ConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: root[:],
	}

	if err := r.BaseRunner.decide(r, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (r *SyncCommitteeRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *SyncCommitteeRunner) GetNetwork() Network {
	return r.network
}

func (r *SyncCommitteeRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *SyncCommitteeRunner) GetShare() *types.Share {
	return r.BaseRunner.Share
}

func (r *SyncCommitteeRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *SyncCommitteeRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
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
func (r *SyncCommitteeRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
