package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type Getters interface {
	GetBaseRunner() *BaseRunner
	GetBeaconNode() BeaconNode
	GetValCheckF() qbft.ProposedValueCheckF
	GetSigner() types.BeaconSigner
	GetOperatorSigner() types.OperatorSigner
	GetNetwork() Network
}

type Runner interface {
	types.Encoder
	types.Root
	Getters

	// StartNewDuty starts a new duty for the runner, returns error if can't
	StartNewDuty(duty types.Duty, quorum uint64) error
	// HasRunningDuty returns true if it has a running duty
	HasRunningDuty() bool
	// ProcessPreConsensus processes all pre-consensus msgs, returns error if can't process
	ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error
	// ProcessConsensus processes all consensus msgs, returns error if can't process
	ProcessConsensus(msg *types.SignedSSVMessage) error
	// ProcessPostConsensus processes all post-consensus msgs, returns error if can't process
	ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error

	// expectedPreConsensusRootsAndDomain an INTERNAL function, returns the expected pre-consensus roots to sign
	expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, spec.DomainType, error)
	// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
	expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, spec.DomainType, error)
	// executeDuty an INTERNAL function, executes a duty.
	executeDuty(duty types.Duty) error
}

type BaseRunner struct {
	State          *State
	Share          map[spec.ValidatorIndex]*types.Share
	QBFTController *qbft.Controller
	BeaconNetwork  types.BeaconNetwork
	RunnerRoleType types.RunnerRole
	types.OperatorSigner

	// highestDecidedSlot holds the highest decided duty slot and gets updated after each decided is reached
	highestDecidedSlot spec.Slot
}

func NewBaseRunner(
	state *State,
	share map[spec.ValidatorIndex]*types.Share,
	controller *qbft.Controller,
	beaconNetwork types.BeaconNetwork,
	runnerRoleType types.RunnerRole,
	operatorSigner types.OperatorSigner,
	highestDecidedSlot spec.Slot,
) *BaseRunner {
	return &BaseRunner{
		State:              state,
		Share:              share,
		QBFTController:     controller,
		BeaconNetwork:      beaconNetwork,
		RunnerRoleType:     runnerRoleType,
		highestDecidedSlot: highestDecidedSlot,
	}
}

// SetHighestDecidedSlot set highestDecidedSlot for base runner
func (b *BaseRunner) SetHighestDecidedSlot(slot spec.Slot) {
	b.highestDecidedSlot = slot
}

// setupForNewDuty is sets the runner for a new duty
func (b *BaseRunner) baseSetupForNewDuty(duty types.Duty, quorum uint64) {
	b.State = NewRunnerState(quorum, duty)
}

// baseStartNewDuty is a base func that all runner implementation can call to start a duty
func (b *BaseRunner) baseStartNewDuty(runner Runner, duty types.Duty, quorum uint64) error {
	if err := b.ShouldProcessDuty(duty); err != nil {
		return errors.Wrap(err, "can't start duty")
	}

	b.baseSetupForNewDuty(duty, quorum)

	return runner.executeDuty(duty)
}

// baseStartNewBeaconDuty is a base func that all runner implementation can call to start a non-beacon duty
func (b *BaseRunner) baseStartNewNonBeaconDuty(runner Runner, duty *types.ValidatorDuty, quorum uint64) error {
	if err := b.ShouldProcessNonBeaconDuty(duty); err != nil {
		return errors.Wrap(err, "can't start non-beacon duty")
	}
	b.baseSetupForNewDuty(duty, quorum)
	return runner.executeDuty(duty)
}

// basePreConsensusMsgProcessing is a base func that all runner implementation can call for processing a pre-consensus msg
func (b *BaseRunner) basePreConsensusMsgProcessing(runner Runner, psigMsgs *types.PartialSignatureMessages) (bool, [][32]byte, error) {
	if err := b.ValidatePreConsensusMsg(runner, psigMsgs); err != nil {
		return false, nil, errors.Wrap(err, "invalid pre-consensus message")
	}

	hasQuorum, roots, err := b.basePartialSigMsgProcessing(psigMsgs, b.State.PreConsensusContainer)
	return hasQuorum, roots, errors.Wrap(err, "could not process pre-consensus partial signature msg")
}

// baseConsensusMsgProcessing is a base func that all runner implementation can call for processing a consensus msg
func (b *BaseRunner) baseConsensusMsgProcessing(runner Runner, msg *types.SignedSSVMessage) (decided bool,
	decidedValue types.Encoder, err error) {
	prevDecided := false
	if b.hasRunningDuty() && b.State != nil && b.State.RunningInstance != nil {
		prevDecided, _ = b.State.RunningInstance.IsDecided()
	}

	decidedSignedMsg, err := b.QBFTController.ProcessMsg(msg)
	if err != nil {
		return false, nil, err
	}

	// we allow all consensus msgs to be processed, once the process finishes we check if there is an actual running duty
	// do not return error if no running duty
	if !b.hasRunningDuty() {
		return false, nil, nil
	}

	if decideCorrectly, err := b.didDecideCorrectly(prevDecided, decidedSignedMsg); !decideCorrectly {
		return false, nil, err
	}

	// decode consensus data
	switch runner.(type) {
	case CommitteeRunner:
		decidedValue = &types.BeaconVote{}
	default:
		decidedValue = &types.ValidatorConsensusData{}
	}
	if err := decidedValue.Decode(decidedSignedMsg.FullData); err != nil {
		return true, nil, errors.Wrap(err, "failed to parse decided value to ValidatorConsensusData")
	}

	if err := b.validateDecidedConsensusData(runner, decidedValue); err != nil {
		return true, nil, errors.Wrap(err, "decided ValidatorConsensusData invalid")
	}

	runner.GetBaseRunner().State.DecidedValue, err = decidedValue.Encode()
	if err != nil {
		return true, nil, errors.Wrap(err, "could not encode decided value")
	}

	// update the highest decided slot
	b.highestDecidedSlot = b.State.StartingDuty.DutySlot()

	return true, decidedValue, nil
}

// basePostConsensusMsgProcessing is a base func that all runner implementation can call for processing a post-consensus msg
// returns whether at least one quorum exists and the roots of the quorums
func (b *BaseRunner) basePostConsensusMsgProcessing(runner Runner, psigMsgs *types.PartialSignatureMessages) (bool, [][32]byte, error) {
	if err := b.ValidatePostConsensusMsg(runner, psigMsgs); err != nil {
		return false, nil, errors.Wrap(err, "invalid post-consensus message")
	}

	hasQuorum, roots, err := b.basePartialSigMsgProcessing(psigMsgs, b.State.PostConsensusContainer)
	return hasQuorum, roots, errors.Wrap(err, "could not process post-consensus partial signature msg")
}

// basePartialSigMsgProcessing adds a validated (without signature verification) partial msg to the container, checks for quorum and returns true (and roots) if quorum exists
func (b *BaseRunner) basePartialSigMsgProcessing(
	psigMsgs *types.PartialSignatureMessages,
	container *PartialSigContainer,
) (bool, [][32]byte, error) {
	roots := make([][32]byte, 0)
	anyQuorum := false
	for _, msg := range psigMsgs.Messages {
		prevQuorum := container.HasQuorum(msg.ValidatorIndex, msg.SigningRoot)

		// Check if it has two signatures for the same signer
		if container.HasSignature(msg.ValidatorIndex, msg.Signer, msg.SigningRoot) {
			b.resolveDuplicateSignature(container, msg)
		} else {
			container.AddSignature(msg)
		}

		hasQuorum := container.HasQuorum(msg.ValidatorIndex, msg.SigningRoot)

		if hasQuorum && !prevQuorum {
			// Notify about first quorum only
			roots = append(roots, msg.SigningRoot)
			anyQuorum = true
		}
	}

	return anyQuorum, roots, nil
}

// didDecideCorrectly returns true if the expected consensus instance decided correctly
func (b *BaseRunner) didDecideCorrectly(prevDecided bool, signedMessage *types.SignedSSVMessage) (bool, error) {
	if signedMessage == nil {
		return false, nil
	}

	if signedMessage.SSVMessage == nil {
		return false, errors.New("ssv message is nil")
	}

	decidedMessage, err := qbft.DecodeMessage(signedMessage.SSVMessage.Data)
	if err != nil {
		return false, err
	}

	if decidedMessage == nil {
		return false, nil
	}

	if b.State.RunningInstance == nil {
		return false, errors.New("decided wrong instance")
	}

	if decidedMessage.Height != b.State.RunningInstance.GetHeight() {
		return false, errors.New("decided wrong instance")
	}

	// verify we decided running instance only, if not we do not proceed
	if prevDecided {
		return false, nil
	}

	return true, nil
}

// decide input param can be a BeaconVote or ValidatorConsensusData
func (b *BaseRunner) decide(runner Runner, slot phase0.Slot, input types.Encoder) error {
	byts, err := input.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode input data for consensus")
	}
	if err := runner.GetValCheckF()(byts); err != nil {
		return errors.Wrap(err, "input data invalid")
	}

	if err := runner.GetBaseRunner().QBFTController.StartNewInstance(
		qbft.Height(slot),
		byts,
	); err != nil {
		return errors.Wrap(err, "could not start new QBFT instance")
	}
	newInstance := runner.GetBaseRunner().QBFTController.InstanceForHeight(runner.GetBaseRunner().QBFTController.Height)
	if newInstance == nil {
		return errors.New("could not find newly created QBFT instance")
	}

	runner.GetBaseRunner().State.RunningInstance = newInstance
	return nil
}

// hasRunningDuty returns true if a new duty didn't start or an existing duty marked as finished
func (b *BaseRunner) hasRunningDuty() bool {
	if b.State == nil {
		return false
	}
	return !b.State.Finished
}

func (b *BaseRunner) ShouldProcessDuty(duty types.Duty) error {
	if b.QBFTController.Height >= qbft.Height(duty.DutySlot()) && b.QBFTController.Height != 0 {
		return errors.Errorf("duty for slot %d already passed. Current height is %d", duty.DutySlot(),
			b.QBFTController.Height)
	}
	return nil
}

func (b *BaseRunner) ShouldProcessNonBeaconDuty(duty types.Duty) error {
	// assume StartingDuty is not nil if state is not nil
	if b.State != nil && b.State.StartingDuty.DutySlot() >= duty.DutySlot() {
		return errors.Errorf("duty for slot %d already passed. Current slot is %d", duty.DutySlot(),
			b.State.StartingDuty.DutySlot())
	}
	return nil
}
