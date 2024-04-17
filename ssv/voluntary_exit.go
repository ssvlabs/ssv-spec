package ssv

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

// BeaconDuty runner for validator voluntary exit duty
type VoluntaryExitRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF

	voluntaryExit *phase0.VoluntaryExit
}

func NewVoluntaryExitRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
) Runner {
	return &VoluntaryExitRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: RoleVoluntaryExit,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
	}
}

func (r *VoluntaryExitRunner) StartNewDuty(duty types.Duty) error {
	// Note: Voluntary exit doesn't require any consensus, it can start a new duty even if previous one didn't finish
	return r.BaseRunner.baseStartNewNonBeaconDuty(r, duty.(*types.BeaconDuty))
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *VoluntaryExitRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

// ProcessPreConsensus Check for quorum of partial signatures over VoluntaryExit and,
// if has quorum, constructs SignedVoluntaryExit and submits to BeaconNode
func (r *VoluntaryExitRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing voluntary exit message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// only 1 root, verified in basePreConsensusMsgProcessing
	root := roots[0]
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root,
		r.GetShare().ValidatorPubKey[:])
	if err != nil {
		// If the reconstructed signature verification failed, fall back to verifying each partial signature
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee)
		return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
	}
	specSig := phase0.BLSSignature{}
	copy(specSig[:], fullSig)

	// create SignedVoluntaryExit using VoluntaryExit created on r.executeDuty() and reconstructed signature
	signedVoluntaryExit := &phase0.SignedVoluntaryExit{
		Message:   r.voluntaryExit,
		Signature: specSig,
	}

	if err := r.beacon.SubmitVoluntaryExit(signedVoluntaryExit); err != nil {
		return errors.Wrap(err, "could not submit voluntary exit")
	}

	r.GetState().Finished = true
	return nil
}

func (r *VoluntaryExitRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return errors.New("no consensus phase for voluntary exit")
}

func (r *VoluntaryExitRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return errors.New("no post consensus phase for voluntary exit")
}

func (r *VoluntaryExitRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	vr, err := r.calculateVoluntaryExit()
	if err != nil {
		return nil, types.DomainError, errors.Wrap(err, "could not calculate voluntary exit")
	}
	return []ssz.HashRoot{vr}, types.DomainVoluntaryExit, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *VoluntaryExitRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, errors.New("no post consensus roots for voluntary exit")
}

// Validator voluntary exit duty doesn't need consensus nor post-consensus.
// It just performs pre-consensus with VoluntaryExitPartialSig over
// a VoluntaryExit object to create a SignedVoluntaryExit
func (r *VoluntaryExitRunner) executeDuty(duty types.Duty) error {
	voluntaryExit, err := r.calculateVoluntaryExit()
	if err != nil {
		return errors.Wrap(err, "could not calculate voluntary exit")
	}

	// get PartialSignatureMessage with voluntaryExit root and signature
	msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.BeaconDuty), voluntaryExit, duty.DutySlot(),
		types.DomainVoluntaryExit)
	if err != nil {
		return errors.Wrap(err, "could not sign VoluntaryExit object")
	}

	msgs := &types.PartialSignatureMessages{
		Type:     types.VoluntaryExitPartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)
	msgToBroadcast, err := types.PartialSignatureMessagesToSignedSSVMessage(msgs, msgID, r.operatorSigner)
	if err != nil {
		return errors.Wrap(err, "could not sign pre-consensus partial signature message")
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast signedPartialMsg with VoluntaryExit")
	}

	// stores value for later using in ProcessPreConsensus
	r.voluntaryExit = voluntaryExit

	return nil
}

// Returns *phase0.VoluntaryExit object with current epoch and own validator index
func (r *VoluntaryExitRunner) calculateVoluntaryExit() (*phase0.VoluntaryExit, error) {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.DutySlot())
	validatorIndex := r.GetState().StartingDuty.(*types.BeaconDuty).ValidatorIndex
	return &phase0.VoluntaryExit{
		Epoch:          epoch,
		ValidatorIndex: validatorIndex,
	}, nil
}

func (r *VoluntaryExitRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *VoluntaryExitRunner) GetNetwork() Network {
	return r.network
}

func (r *VoluntaryExitRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *VoluntaryExitRunner) GetShare() *types.Share {
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *VoluntaryExitRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *VoluntaryExitRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *VoluntaryExitRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *VoluntaryExitRunner) GetOperatorSigner() types.OperatorSigner {
	return r.operatorSigner
}

// Encode returns the encoded struct in bytes or error
func (r *VoluntaryExitRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *VoluntaryExitRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *VoluntaryExitRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
