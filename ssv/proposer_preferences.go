package ssv

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ProposerPreferencesRunner runs the Gloas (ePBS) proposer-preferences duty (SIP #94 §5). A validator
// legitimately holds several upcoming proposal slots in the lookahead at once, and preference messages
// route by MessageID (validator + role) with the proposal slot carried inside the message — so unlike
// the single-duty runners this one dispatches each proposal slot to its own
// ProposerPreferencesSlotRunner; a single-state runner would make concurrent lookahead slots overwrite
// or reject one another.
//
// The embedded BaseRunner provides only the static runner surface (role, share); the per-slot
// sub-runners own the duty state.
type ProposerPreferencesRunner struct {
	BaseRunner *BaseRunner

	// BySlot holds one sub-runner per concurrently-active proposal slot.
	BySlot map[phase0.Slot]*ProposerPreferencesSlotRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner

	gasLimit uint64
}

func NewProposerPreferencesRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	gasLimit uint64,
) (Runner, error) {

	if len(share) != 1 {
		return nil, fmt.Errorf("must have one share")
	}
	if err := validateShareMap(share); err != nil {
		return nil, err
	}

	return &ProposerPreferencesRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RoleProposerPreferences,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},
		BySlot: map[phase0.Slot]*ProposerPreferencesSlotRunner{},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		gasLimit:       gasLimit,
	}, nil
}

// StartNewDuty starts an independent per-slot flow for the duty's proposal slot. A duty for a slot
// that already has one is a re-emission (e.g. after a reorg moved the duty's dependent root): the
// replacement freezes a freshly derived preference and starts a fresh signature container, discarding
// the prior incarnation's.
func (r *ProposerPreferencesRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	sub := r.NewSlotRunner()
	if err := sub.StartNewDuty(duty, quorum); err != nil {
		return err
	}
	r.BySlot[duty.DutySlot()] = sub
	return nil
}

// NewSlotRunner builds a sub-runner sharing the dispatcher's dependencies; StartNewDuty creates one
// per started duty, and the spec-test harness seeds pre-existing slot flows through it.
func (r *ProposerPreferencesRunner) NewSlotRunner() *ProposerPreferencesSlotRunner {
	return NewProposerPreferencesSlotRunner(r.BaseRunner.BeaconNetwork, r.BaseRunner.Share, r.beacon,
		r.network, r.signer, r.operatorSigner, r.gasLimit)
}

// HasRunningDuty returns true if any proposal slot's duty is still running
func (r *ProposerPreferencesRunner) HasRunningDuty() bool {
	for _, sub := range r.BySlot {
		if sub.HasRunningDuty() {
			return true
		}
	}
	return false
}

// ProcessPreConsensus routes the message to its proposal slot's sub-runner
func (r *ProposerPreferencesRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	sub, ok := r.BySlot[signedMsg.Slot]
	if !ok {
		return types.NewError(types.NoRunningDutyErrorCode, "no duty for the message's proposal slot")
	}
	return sub.ProcessPreConsensus(signedMsg)
}

func (r *ProposerPreferencesRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return types.NewError(types.ProposerPreferencesNoConsensusPhaseErrorCode, "no consensus phase for proposer preferences")
}

func (r *ProposerPreferencesRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.ProposerPreferencesNoPostConsensusPhaseErrorCode, "no post consensus phase for proposer preferences")
}

// expectedPreConsensusRootsAndDomain / expectedPostConsensusRootsAndDomain / executeDuty run on the
// per-slot sub-runners, never on the dispatcher.
func (r *ProposerPreferencesRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, types.DomainError, fmt.Errorf("proposer preferences dispatcher has no frozen preference")
}

func (r *ProposerPreferencesRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, fmt.Errorf("no post consensus roots for proposer preferences")
}

func (r *ProposerPreferencesRunner) executeDuty(duty types.Duty) error {
	return fmt.Errorf("proposer preferences dispatcher does not execute duties directly")
}

func (r *ProposerPreferencesRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *ProposerPreferencesRunner) GetNetwork() Network {
	return r.network
}

func (r *ProposerPreferencesRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *ProposerPreferencesRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *ProposerPreferencesRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *ProposerPreferencesRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return nil
}

func (r *ProposerPreferencesRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *ProposerPreferencesRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}

// ProposerPreferencesSlotRunner is one proposal slot's independent preference flow: derive and freeze
// the preference, sign it under DomainProposerPreferences, and submit once a quorum of operators
// converged on the same preference (honest convergence over the shared config and observed dependent
// root — divergence splits signing roots and costs liveness, never safety).
type ProposerPreferencesSlotRunner struct {
	BaseRunner *BaseRunner

	// ProposerPreferences is this slot's frozen preference. Incoming partial signatures are validated
	// against exactly its root; nil until the duty executes.
	ProposerPreferences *gloas.ProposerPreferences

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner

	gasLimit uint64
}

// NewProposerPreferencesSlotRunner constructs one proposal slot's flow; the dispatcher creates one
// per started duty (exported for the spec-test harness, which rebuilds serialized sub-runners).
func NewProposerPreferencesSlotRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	gasLimit uint64,
) *ProposerPreferencesSlotRunner {
	return &ProposerPreferencesSlotRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RoleProposerPreferences,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		gasLimit:       gasLimit,
	}
}

func (r *ProposerPreferencesSlotRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	return r.BaseRunner.baseStartNewNonBeaconDuty(r, duty.(*types.ValidatorDuty), quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *ProposerPreferencesSlotRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *ProposerPreferencesSlotRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing proposer preferences message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// Defensive: peer partials are validated against the frozen preference, so a quorum without one
	// should be unreachable.
	if r.ProposerPreferences == nil {
		return types.NewError(types.ProposerPreferencesNoPreferenceErrorCode, "reached quorum without a frozen preference")
	}

	// only 1 root, verified in basePreConsensusMsgProcessing
	root := roots[0]
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
	if err != nil {
		// If the reconstructed signature verification failed, fall back to verifying each partial signature
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee,
			r.GetShare().ValidatorIndex)
		return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
	}
	specSig := phase0.BLSSignature{}
	copy(specSig[:], fullSig)

	signed := &gloas.SignedProposerPreferences{
		Message:   r.ProposerPreferences,
		Signature: specSig,
	}

	if err := r.beacon.SubmitProposerPreferences(signed); err != nil {
		return errors.Wrap(err, "could not submit proposer preferences")
	}

	r.GetState().Finished = true
	return nil
}

func (r *ProposerPreferencesSlotRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return types.NewError(types.ProposerPreferencesNoConsensusPhaseErrorCode, "no consensus phase for proposer preferences")
}

func (r *ProposerPreferencesSlotRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.ProposerPreferencesNoPostConsensusPhaseErrorCode, "no post consensus phase for proposer preferences")
}

func (r *ProposerPreferencesSlotRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	if r.ProposerPreferences == nil {
		return nil, types.DomainError, types.NewError(types.ProposerPreferencesNoPreferenceErrorCode, "no frozen preference")
	}
	return []ssz.HashRoot{r.ProposerPreferences}, types.DomainProposerPreferences, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ProposerPreferencesSlotRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, fmt.Errorf("no post consensus roots for proposer preferences")
}

// executeDuty steps:
//  1. derive and freeze the slot's preference: the proposer-duties dependent root from the local
//     beacon node, the share's fee recipient, and the configured target gas limit (SIP #94 §5)
//  2. sign it under DomainProposerPreferences — the domain epoch is the proposal slot's epoch even
//     when emitted earlier, which is what makes pre-fork emission for post-fork slots work — and
//     broadcast the partial signature with the proposal slot
//  3. once a quorum of operators converged on the same preference, reconstruct and submit it
func (r *ProposerPreferencesSlotRunner) executeDuty(duty types.Duty) error {
	proposalSlot := duty.DutySlot()
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(proposalSlot)

	dependentRoot, err := r.beacon.ProposerDutiesDependentRoot(epoch)
	if err != nil {
		return errors.Wrap(err, "failed to get proposer duties dependent root")
	}

	preferences := &gloas.ProposerPreferences{
		DependentRoot:  dependentRoot,
		ProposalSlot:   proposalSlot,
		ValidatorIndex: r.GetShare().ValidatorIndex,
		FeeRecipient:   r.GetShare().FeeRecipientAddress,
		TargetGasLimit: r.gasLimit,
	}
	r.ProposerPreferences = preferences

	msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.ValidatorDuty), preferences, proposalSlot,
		types.DomainProposerPreferences)
	if err != nil {
		return errors.Wrap(err, "could not sign proposer preferences")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.ProposerPreferencesPartialSig,
		Slot:     proposalSlot,
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewValidatorMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.RunnerRoleType)

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
		return errors.Wrap(err, "can't broadcast partial proposer preferences sig")
	}
	return nil
}

func (r *ProposerPreferencesSlotRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *ProposerPreferencesSlotRunner) GetNetwork() Network {
	return r.network
}

func (r *ProposerPreferencesSlotRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *ProposerPreferencesSlotRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *ProposerPreferencesSlotRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *ProposerPreferencesSlotRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return nil
}

func (r *ProposerPreferencesSlotRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *ProposerPreferencesSlotRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}
