package ssv

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
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

	// BuilderEntries is the cluster-configured §5 builder-request-auth entry set (static config; SIP #94
	// §5). It is exported so it survives the spec-test JSON round-trip: unlike gasLimit (a fixed test
	// constant reconstructed by role), it varies per test, and a reconstructed dispatcher must carry it to
	// re-freeze the auth round. The dispatcher passes it to each sub-runner it creates.
	BuilderEntries []BuilderEntry

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
	builderEntries []BuilderEntry,
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
		BySlot:         map[phase0.Slot]*ProposerPreferencesSlotRunner{},
		BuilderEntries: builderEntries,

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
		r.network, r.signer, r.operatorSigner, r.gasLimit, r.BuilderEntries)
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

// ProcessEnvelopeDissemination returns an error: only the envelope-proposer runner processes
// disseminated envelopes (SIP #94 §6).
func (*ProposerPreferencesRunner) ProcessEnvelopeDissemination(*types.SignedSSVMessage) error {
	return types.NewError(types.EnvelopeDisseminationUnsupportedErrorCode, "runner does not process envelope dissemination")
}

func (r *ProposerPreferencesRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.ProposerPreferencesNoPostConsensusPhaseErrorCode, "no post consensus phase for proposer preferences")
}

// expectedPreConsensusRootsAndDomain / expectedPostConsensusRootsAndDomain / executeDuty run on the
// per-slot sub-runners, never on the dispatcher.
func (r *ProposerPreferencesRunner) expectedPreConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	return nil, types.DomainError, fmt.Errorf("proposer preferences dispatcher has no frozen preference")
}

func (r *ProposerPreferencesRunner) expectedPostConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
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

	// ProposerPreferences is this slot's frozen preference. Incoming ProposerPreferencesPartialSig
	// partials are validated against exactly its root; nil until the duty executes.
	ProposerPreferences *gloas.ProposerPreferences
	// BuilderRequestAuths is this slot's frozen builder-request-auth set — one per distinct configured
	// entry data (SIP #94 §5 builder-request-auth extension). Incoming RequestAuthPartialSig partials
	// validate against these roots and collect per root, independently of the preference round. Empty
	// when no builder entries are configured.
	BuilderRequestAuths []*gloas.BuilderRequestAuth

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner

	gasLimit       uint64
	builderEntries []BuilderEntry
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
	builderEntries []BuilderEntry,
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
		builderEntries: builderEntries,
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
	// The builder-request-auth round rides the same duty but collects independently of the preference
	// round — neither gates the other (SIP #94 §5).
	if signedMsg.Type == types.RequestAuthPartialSig {
		return r.processRequestAuth(signedMsg)
	}

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

// ProcessEnvelopeDissemination returns an error: only the envelope-proposer runner processes
// disseminated envelopes (SIP #94 §6).
func (*ProposerPreferencesSlotRunner) ProcessEnvelopeDissemination(*types.SignedSSVMessage) error {
	return types.NewError(types.EnvelopeDisseminationUnsupportedErrorCode, "runner does not process envelope dissemination")
}

func (r *ProposerPreferencesSlotRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.ProposerPreferencesNoPostConsensusPhaseErrorCode, "no post consensus phase for proposer preferences")
}

func (r *ProposerPreferencesSlotRunner) expectedPreConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	if r.ProposerPreferences == nil {
		return nil, types.DomainError, types.NewError(types.ProposerPreferencesNoPreferenceErrorCode, "no frozen preference")
	}
	return []types.HashRoot{r.ProposerPreferences}, types.DomainProposerPreferences, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ProposerPreferencesSlotRunner) expectedPostConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, fmt.Errorf("no post consensus roots for proposer preferences")
}

// executeDuty steps:
//  1. derive and freeze the slot's preference: the proposer-duties dependent root from the local
//     beacon node, the share's fee recipient, and the configured target gas limit (SIP #94 §5)
//  2. sign it under DomainProposerPreferences — the domain epoch is the proposal slot's epoch even
//     when emitted earlier, which is what makes pre-fork emission for post-fork slots work — and
//     broadcast the partial signature with the proposal slot
//  3. once a quorum of operators converged on the same preference, reconstruct and submit it
//  4. run the builder-request-auth round riding this duty (executeRequestAuthRound), independently of
//     the preference round
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
	if err := r.broadcastPartialSig(msgs); err != nil {
		return err
	}

	return r.executeRequestAuthRound(duty.(*types.ValidatorDuty), proposalSlot)
}

// executeRequestAuthRound freezes one BuilderRequestAuth per distinct configured entry data, signs each
// under DomainBuilderRequestAuth (chain-independent), and broadcasts them together in a single
// RequestAuthPartialSig container — one partial per frozen auth, the multi-root pre-consensus shape the
// container and quorum machinery expect. Entries sharing data share a root; zero-length data is skipped;
// entries are capped at MaxBuilderEntries (SIP #94 §5). No configured entries means no auth round.
func (r *ProposerPreferencesSlotRunner) executeRequestAuthRound(duty *types.ValidatorDuty, proposalSlot phase0.Slot) error {
	r.BuilderRequestAuths = nil
	authMsgs := &types.PartialSignatureMessages{
		Type:     types.RequestAuthPartialSig,
		Slot:     proposalSlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	seen := make(map[string]bool)
	for _, entry := range r.builderEntries {
		if len(r.BuilderRequestAuths) >= MaxBuilderEntries {
			break
		}
		data := entry.AuthData()
		if len(data) == 0 || seen[string(data)] {
			continue
		}
		seen[string(data)] = true

		auth := &gloas.BuilderRequestAuth{Data: data, Slot: proposalSlot}
		r.BuilderRequestAuths = append(r.BuilderRequestAuths, auth)

		msg, err := r.BaseRunner.signBeaconObject(r, duty, auth, proposalSlot, types.DomainBuilderRequestAuth)
		if err != nil {
			return errors.Wrap(err, "could not sign builder request auth")
		}
		authMsgs.Messages = append(authMsgs.Messages, msg)
	}
	if len(authMsgs.Messages) == 0 {
		return nil
	}
	return r.broadcastPartialSig(authMsgs)
}

// broadcastPartialSig operator-signs a partial-signature container into an SSVMessage and broadcasts it.
func (r *ProposerPreferencesSlotRunner) broadcastPartialSig(msgs *types.PartialSignatureMessages) error {
	encodedMsg, err := msgs.Encode()
	if err != nil {
		return err
	}
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewValidatorMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.RunnerRoleType),
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
		return errors.Wrap(err, "can't broadcast partial signature")
	}
	return nil
}

// processRequestAuth handles a RequestAuthPartialSig container (SIP #94 §5): it collects its partials
// against the frozen auth roots — independently of the preference round and of the runner's finished
// state, since auth collection continues until the proposal slot — and reconstructs and submits one
// SignedBuilderRequestAuth per root that reaches quorum.
func (r *ProposerPreferencesSlotRunner) processRequestAuth(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.baseRequestAuthProcessing(signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing builder request auth message")
	}
	if !quorum {
		return nil
	}

	// A multi-root container can push several auth roots over quorum at once; submit each.
	for _, root := range roots {
		auth := r.authForSigningRoot(root)
		if auth == nil {
			return types.NewError(types.RequestAuthNoAuthErrorCode, "builder-request-auth quorum for an unknown root")
		}

		fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
		if err != nil {
			r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee,
				r.GetShare().ValidatorIndex)
			return errors.Wrap(err, "got builder-request-auth quorum but it has invalid signatures")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], fullSig)

		if err := r.beacon.SubmitBuilderRequestAuth(&gloas.SignedBuilderRequestAuth{Message: auth, Signature: specSig}); err != nil {
			return errors.Wrap(err, "could not submit builder request auth")
		}
	}
	return nil
}

// baseRequestAuthProcessing validates a RequestAuthPartialSig against the frozen auth roots under
// DomainBuilderRequestAuth and adds it to the pre-consensus container (per-root quorum). It deliberately
// skips the running-duty check the preference round uses: the auth round keeps collecting after the
// preference round finishes (SIP #94 §5 — neither gates the other).
func (r *ProposerPreferencesSlotRunner) baseRequestAuthProcessing(signedMsg *types.PartialSignatureMessages) (bool, [][32]byte, error) {
	if r.BaseRunner.State == nil {
		return false, nil, types.NewError(types.NoRunningDutyErrorCode, "no running duty")
	}
	if err := r.BaseRunner.validatePartialSigMsgForSlot(signedMsg, r.BaseRunner.State.StartingDuty.DutySlot()); err != nil {
		return false, nil, err
	}
	if err := r.BaseRunner.validateValidatorIndexInPartialSigMsg(signedMsg); err != nil {
		return false, nil, err
	}
	roots, domain, err := r.expectedRequestAuthRootsAndDomain()
	if err != nil {
		return false, nil, err
	}
	if err := r.BaseRunner.verifyExpectedRoot(r, signedMsg, roots, domain); err != nil {
		return false, nil, err
	}
	quorum, quorumRoots := r.BaseRunner.basePartialSigMsgProcessing(signedMsg, r.GetState().PreConsensusContainer)
	return quorum, quorumRoots, nil
}

// expectedRequestAuthRootsAndDomain returns the frozen builder-request-auth roots and their domain, so
// incoming RequestAuthPartialSig partials validate against exactly this operator's configured entries.
func (r *ProposerPreferencesSlotRunner) expectedRequestAuthRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	if len(r.BuilderRequestAuths) == 0 {
		return nil, types.DomainError, types.NewError(types.RequestAuthNoAuthErrorCode, "no frozen builder request auths")
	}
	roots := make([]types.HashRoot, 0, len(r.BuilderRequestAuths))
	for _, auth := range r.BuilderRequestAuths {
		roots = append(roots, auth)
	}
	return roots, types.DomainBuilderRequestAuth, nil
}

// authForSigningRoot returns the frozen auth whose DomainBuilderRequestAuth signing root equals root.
func (r *ProposerPreferencesSlotRunner) authForSigningRoot(root [32]byte) *gloas.BuilderRequestAuth {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.DutySlot())
	domain, err := r.beacon.DomainData(epoch, types.DomainBuilderRequestAuth)
	if err != nil {
		return nil
	}
	for _, auth := range r.BuilderRequestAuths {
		signingRoot, err := types.ComputeETHSigningRoot(auth, domain)
		if err != nil {
			continue
		}
		if signingRoot == root {
			return auth
		}
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
