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

	// BySlot holds one sub-runner per started proposal slot. Entries are never removed: pruning them, and
	// closing a slot's §5/§7 acceptance window at its late end, is node-side lifecycle.
	BySlot map[phase0.Slot]*ProposerPreferencesSlotRunner

	// PendingPartials keeps every partial received, per proposal slot, as single-entry messages that
	// StartNewDuty replays into the slot's new sub-runner. A peer's partial is delivered once (§7 dedups
	// repeats), so this is how a late-starting duty gets the partials that arrived before it (SIP #94 §5)
	// and a re-emission gets those its predecessor rejected under a since-changed preference. Entries are
	// deduplicated by (signer, signing root) and capped at maxPendingPartialsPerSigner per signer and slot;
	// the slot is the sender's choice, so pruning past slots (node-side, like BySlot) is what bounds the
	// map. Omitted from JSON when empty, so runners holding none keep their state roots.
	PendingPartials map[phase0.Slot][]*types.PartialSignatureMessages `json:"PendingPartials,omitempty"`

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

// MaxProposerPreferencesDistinctRoots is SIP #94 §7's per-signer budget of distinct preference roots per
// proposal slot; extra roots come from re-emissions after a preference-input change.
const MaxProposerPreferencesDistinctRoots = 4

// maxPendingPartialsPerSigner caps a signer's PendingPartials per slot at the §7 per-signer budgets combined:
// MaxProposerPreferencesDistinctRoots preference roots plus MaxBuilderEntries request-auth roots. §7 message
// validation enforces the two budgets separately.
const maxPendingPartialsPerSigner = MaxProposerPreferencesDistinctRoots + MaxBuilderEntries

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
		BySlot:          map[phase0.Slot]*ProposerPreferencesSlotRunner{},
		PendingPartials: map[phase0.Slot][]*types.PartialSignatureMessages{},
		BuilderEntries:  builderEntries,

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		gasLimit:       gasLimit,
	}, nil
}

// StartNewDuty runs the duty's proposal slot in a fresh sub-runner, then replays the slot's
// PendingPartials into it.
//
// A duty for a slot that already has a sub-runner is a re-emission (e.g. after a reorg). The fresh sub
// re-derives its roots and carries the prior sub's collected shares over per signing root, so a quorum
// already reached does not fire (and submit) again. Auth roots carry no dependent_root, so their shares
// always carry; preference shares carry only if the re-derived preference is unchanged.
//
// Left to node-side retry: a root the prior sub had at quorum but failed to submit is not re-submitted
// (quorum fires once), and a re-emission whose preference derivation fails has no preference round until
// a later re-emission.
func (r *ProposerPreferencesRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	slot := duty.DutySlot()
	prev := r.BySlot[slot]
	sub := r.NewSlotRunner()
	// Install the sub before executing, so it collects peers' shares even if executing the duty fails.
	r.BySlot[slot] = sub
	startErr := sub.StartNewDuty(duty, quorum)
	// Carry over even on failure: the two rounds are independent (SIP #94 §5), so one failing must not
	// discard the other's shares.
	if prev != nil {
		sub.carryOverAuthShares(prev)
		sub.carryOverPreferenceShares(prev)
	}
	// Replay each retained partial as if it just arrived. One that fails its checks (e.g. a diverging
	// peer's) is skipped. A failure to reconstruct or submit a quorum it completes is returned, as on the
	// live path, after the duty's own error.
	var replayErr error
	for _, msg := range r.PendingPartials[slot] {
		roots, err := sub.collectPartials(msg)
		if err != nil {
			continue
		}
		if err := sub.submitQuorumRoots(msg.Type, roots); err != nil && replayErr == nil {
			replayErr = err
		}
	}
	if startErr != nil {
		return startErr
	}
	return replayErr
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

// ProcessPreConsensus retains the message's partials for its proposal slot (PendingPartials), then routes
// it to the slot's sub-runner if that slot's duty has started here.
func (r *ProposerPreferencesRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	// Reject any type other than the duty's two rounds before retaining it, so it takes no room.
	if signedMsg.Type != types.ProposerPreferencesPartialSig && signedMsg.Type != types.RequestAuthPartialSig {
		return types.NewError(types.ProposerPreferencesUnexpectedPartialSigTypeErrorCode,
			"unexpected pre-consensus partial signature type for proposer-preferences runner")
	}
	r.retainPartials(signedMsg)

	sub, ok := r.BySlot[signedMsg.Slot]
	if !ok {
		// Its duty has not started here yet; StartNewDuty replays the retained partials once it does.
		return nil
	}
	return sub.ProcessPreConsensus(signedMsg)
}

// retainPartials adds each of the message's partials to its slot's PendingPartials as a single-entry
// message, skipping a (signer, signing root) already retained and any beyond the signer's cap. The two
// types sign under different domains, so one keyspace serves both.
func (r *ProposerPreferencesRunner) retainPartials(signedMsg *types.PartialSignatureMessages) {
	for _, msg := range signedMsg.Messages {
		r.retainPartial(signedMsg, msg)
	}
}

// retainPartial adds one partial of signedMsg to PendingPartials (see retainPartials). The cap counts per
// signer, so a faulty operator cannot crowd out the others' partials.
func (r *ProposerPreferencesRunner) retainPartial(signedMsg *types.PartialSignatureMessages, msg *types.PartialSignatureMessage) {
	retained := r.PendingPartials[signedMsg.Slot]
	perSigner := 0
	for _, existing := range retained {
		if existing.Messages[0].Signer != msg.Signer {
			continue
		}
		if existing.Messages[0].SigningRoot == msg.SigningRoot {
			return
		}
		perSigner++
	}
	if perSigner >= maxPendingPartialsPerSigner {
		return
	}
	r.PendingPartials[signedMsg.Slot] = append(retained, &types.PartialSignatureMessages{
		Type:     signedMsg.Type,
		Slot:     signedMsg.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
	})
}

func (r *ProposerPreferencesRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return types.NewError(types.ProposerPreferencesNoConsensusPhaseErrorCode, "no consensus phase for proposer preferences")
}

func (r *ProposerPreferencesRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.ProposerPreferencesNoPostConsensusPhaseErrorCode, "no post consensus phase for proposer preferences")
}

// expectedPreConsensusRootsAndDomain / expectedPostConsensusRootsAndDomains / executeDuty run on the
// per-slot sub-runners, never on the dispatcher.
func (r *ProposerPreferencesRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, types.DomainError, fmt.Errorf("proposer preferences dispatcher has no frozen preference")
}

func (r *ProposerPreferencesRunner) expectedPostConsensusRootsAndDomains() ([]PostConsensusRoot, error) {
	return nil, fmt.Errorf("no post consensus roots for proposer preferences")
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

// ProcessPreConsensus collects the message's partials, then reconstructs and submits whatever they bring
// to quorum.
func (r *ProposerPreferencesSlotRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	roots, err := r.collectPartials(signedMsg)
	if err != nil {
		return err
	}
	return r.submitQuorumRoots(signedMsg.Type, roots)
}

// collectPartials validates the message and adds its partials to the pre-consensus container, returning
// the roots that just reached quorum (each reported once, on its first crossing). It routes by type between
// the duty's two rounds (SIP #94 §5) and rejects any other type rather than treat it as a preference partial.
func (r *ProposerPreferencesSlotRunner) collectPartials(signedMsg *types.PartialSignatureMessages) ([][32]byte, error) {
	var quorum bool
	var roots [][32]byte
	var err error
	switch signedMsg.Type {
	case types.ProposerPreferencesPartialSig:
		if quorum, roots, err = r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg); err != nil {
			return nil, errors.Wrap(err, "failed processing proposer preferences message")
		}
	case types.RequestAuthPartialSig:
		if quorum, roots, err = r.baseRequestAuthProcessing(signedMsg); err != nil {
			return nil, errors.Wrap(err, "failed processing builder request auth message")
		}
	default:
		return nil, types.NewError(types.ProposerPreferencesUnexpectedPartialSigTypeErrorCode,
			"unexpected pre-consensus partial signature type for proposer-preferences runner")
	}
	if !quorum {
		return nil, nil
	}
	return roots, nil
}

// submitQuorumRoots reconstructs and submits what the roots that just reached quorum sign: the preference,
// or a SignedBuilderRequestAuth per auth root.
func (r *ProposerPreferencesSlotRunner) submitQuorumRoots(msgType types.PartialSigMsgType, roots [][32]byte) error {
	if len(roots) == 0 {
		return nil
	}
	if msgType == types.RequestAuthPartialSig {
		return r.submitRequestAuths(roots)
	}
	// only 1 root, verified in basePreConsensusMsgProcessing
	return r.submitPreferences(roots[0])
}

// submitPreferences reconstructs the signature over the frozen ProposerPreferences once its root reached
// quorum, and submits the SignedProposerPreferences (SIP #94 §5).
func (r *ProposerPreferencesSlotRunner) submitPreferences(root [32]byte) error {
	// Defensive: peer partials are validated against the frozen preference, so a quorum without one
	// should be unreachable.
	if r.ProposerPreferences == nil {
		return types.NewError(types.ProposerPreferencesNoPreferenceErrorCode, "reached quorum without a frozen preference")
	}

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

// expectedPostConsensusRootsAndDomains an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ProposerPreferencesSlotRunner) expectedPostConsensusRootsAndDomains() ([]PostConsensusRoot, error) {
	return nil, fmt.Errorf("no post consensus roots for proposer preferences")
}

// executeDuty runs the two independent rounds this duty carries (SIP #94 §5 — neither gates the other):
// the proposer-preferences round (executePreferenceRound) and the builder-request-auth round
// (executeRequestAuthRound). The auth round runs even when the preference round fails — a dependent-root
// fetch or preference-broadcast failure must not strand the auth round, whose roots carry no dependent_root.
// The preference error is surfaced first, as it is the duty's primary output.
func (r *ProposerPreferencesSlotRunner) executeDuty(duty types.Duty) error {
	vDuty := duty.(*types.ValidatorDuty)
	proposalSlot := duty.DutySlot()
	prefErr := r.executePreferenceRound(vDuty, proposalSlot)
	authErr := r.executeRequestAuthRound(vDuty, proposalSlot)
	if prefErr != nil {
		return prefErr
	}
	return authErr
}

// executePreferenceRound derives, signs and broadcasts the slot's ProposerPreferences (SIP #94 §5):
//  1. freeze the preference: the proposer-duties dependent root from the local beacon node, the share's
//     fee recipient, and the configured target gas limit
//  2. sign it under DomainProposerPreferences — the domain epoch is the proposal slot's epoch even when
//     emitted earlier, which is what makes pre-fork emission for post-fork slots work — and broadcast the
//     partial signature with the proposal slot
//  3. once a quorum of operators converged on the same preference, reconstruct and submit it
func (r *ProposerPreferencesSlotRunner) executePreferenceRound(duty *types.ValidatorDuty, proposalSlot phase0.Slot) error {
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

	msg, err := r.BaseRunner.signBeaconObject(r, duty, preferences, proposalSlot, types.DomainProposerPreferences)
	if err != nil {
		return errors.Wrap(err, "could not sign proposer preferences")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.ProposerPreferencesPartialSig,
		Slot:     proposalSlot,
		Messages: []*types.PartialSignatureMessage{msg},
	}
	return r.broadcastPartialSig(msgs)
}

// executeRequestAuthRound freezes one BuilderRequestAuth per distinct configured entry data, signs each
// under DomainBuilderRequestAuth (chain-independent), and broadcasts them together in a single multi-entry
// RequestAuthPartialSig container (SIP #94 §5/§7: one bounded packet per slot, up to MaxBuilderEntries;
// per-builder isolation is enforced on the receive side per entry). Entries sharing data share a root;
// zero-length data is skipped; entries are capped at MaxBuilderEntries. The round is
// best-effort and never fails the independent preference round (already broadcast by the caller): a
// malformed entry that fails to sign (e.g. oversized Data) is skipped, and an auth is frozen only once its
// share is in the outgoing container, so peer shares are never expected for an entry this operator dropped.
// No configured entries, or none that sign, means no auth broadcast.
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
		msg, err := r.BaseRunner.signBeaconObject(r, duty, auth, proposalSlot, types.DomainBuilderRequestAuth)
		if err != nil {
			continue
		}
		authMsgs.Messages = append(authMsgs.Messages, msg)
		r.BuilderRequestAuths = append(r.BuilderRequestAuths, auth)
	}
	if len(authMsgs.Messages) == 0 {
		return nil
	}
	// Best-effort: a broadcast failure must not fail the independent preference round.
	_ = r.broadcastPartialSig(authMsgs)
	return nil
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

// submitRequestAuths reconstructs and submits a SignedBuilderRequestAuth per auth root that reached quorum
// (SIP #94 §5).
func (r *ProposerPreferencesSlotRunner) submitRequestAuths(roots [][32]byte) error {
	// One packet can push several auth roots over quorum, and quorum fires once per root, so a bad share or
	// a submit failure on one root must not strand the others: attempt every root, falling back to
	// verify-each (which drops invalid shares) when one fails to reconstruct, and return the first error.
	var firstErr error
	record := func(err error) {
		if firstErr == nil {
			firstErr = err
		}
	}
	for _, root := range roots {
		auth := r.authForSigningRoot(root)
		if auth == nil {
			record(types.NewError(types.RequestAuthNoAuthErrorCode, "builder-request-auth quorum for an unknown root"))
			continue
		}

		fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
		if err != nil {
			r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee,
				r.GetShare().ValidatorIndex)
			record(errors.Wrap(err, "got builder-request-auth quorum but it has invalid signatures"))
			continue
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], fullSig)

		if err := r.beacon.SubmitBuilderRequestAuth(&gloas.SignedBuilderRequestAuth{Message: auth, Signature: specSig}); err != nil {
			record(errors.Wrap(err, "could not submit builder request auth"))
			continue
		}
	}
	return firstErr
}

// baseRequestAuthProcessing validates a RequestAuthPartialSig against this operator's frozen auth roots
// under DomainBuilderRequestAuth and adds the matching entries to the pre-consensus container (per-root
// quorum). It skips the running-duty check the preference round uses, so auth collection continues after
// the preference round finishes (SIP #94 §5: neither gates the other). Closing collection at the proposal
// slot (§5) and the 2-slot lateness TTL (§7) are node-side lifecycle, as for BySlot.
//
// A packet carries up to MaxBuilderEntries entries (SIP #94 §5/§7). Entries whose root is not a frozen
// auth are ignored rather than failing the packet, so a divergent entry costs only that builder. The
// shared verifyExpectedRoot does not fit: it requires the full expected set, rejecting a peer that carries
// a subset.
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
	if len(r.BuilderRequestAuths) == 0 {
		return false, nil, types.NewError(types.RequestAuthNoAuthErrorCode, "no frozen builder request auths")
	}
	if len(signedMsg.Messages) > MaxBuilderEntries {
		return false, nil, types.NewError(types.RequestAuthWrongRootsCountErrorCode, "builder-request-auth container exceeds MaxBuilderEntries")
	}
	matched := &types.PartialSignatureMessages{
		Type:     signedMsg.Type,
		Slot:     signedMsg.Slot,
		Messages: make([]*types.PartialSignatureMessage, 0, len(signedMsg.Messages)),
	}
	for _, m := range signedMsg.Messages {
		if r.authForSigningRoot(m.SigningRoot) != nil {
			matched.Messages = append(matched.Messages, m)
		}
	}
	if len(matched.Messages) == 0 {
		return false, nil, types.NewError(types.WrongSigningRootErrorCode, "no builder-request-auth entry matches a frozen auth root")
	}
	quorum, quorumRoots := r.BaseRunner.basePartialSigMsgProcessing(matched, r.GetState().PreConsensusContainer)
	return quorum, quorumRoots, nil
}

// authSigningRoot computes an auth's DomainBuilderRequestAuth signing root for the running duty's epoch.
func (r *ProposerPreferencesSlotRunner) authSigningRoot(auth *gloas.BuilderRequestAuth) ([32]byte, error) {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.DutySlot())
	domain, err := r.beacon.DomainData(epoch, types.DomainBuilderRequestAuth)
	if err != nil {
		return [32]byte{}, err
	}
	return types.ComputeETHSigningRoot(auth, domain)
}

// authForSigningRoot returns the frozen auth whose DomainBuilderRequestAuth signing root equals root.
func (r *ProposerPreferencesSlotRunner) authForSigningRoot(root [32]byte) *gloas.BuilderRequestAuth {
	for _, auth := range r.BuilderRequestAuths {
		signingRoot, err := r.authSigningRoot(auth)
		if err != nil {
			continue
		}
		if signingRoot == root {
			return auth
		}
	}
	return nil
}

// carryOverAuthShares copies already-collected builder-request-auth partial signatures from a prior
// same-slot sub-runner (a re-emission) into this one. Auth roots carry no dependent_root (SIP #94 §5), so
// a re-emission re-freezes byte-identical roots, and a root already at quorum stays there without firing
// (and submitting) again when the retained partials are replayed.
func (r *ProposerPreferencesSlotRunner) carryOverAuthShares(prev *ProposerPreferencesSlotRunner) {
	if prev == nil || prev.GetState() == nil || r.GetState() == nil {
		return
	}
	for _, auth := range r.BuilderRequestAuths {
		signingRoot, err := r.authSigningRoot(auth)
		if err != nil {
			continue
		}
		r.carryOverSharesForRoot(prev, signingRoot)
	}
}

// preferenceSigningRoot is the signing root of this sub's frozen ProposerPreferences under
// DomainProposerPreferences — the root peers sign in the preference round.
func (r *ProposerPreferencesSlotRunner) preferenceSigningRoot() ([32]byte, error) {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.DutySlot())
	domain, err := r.beacon.DomainData(epoch, types.DomainProposerPreferences)
	if err != nil {
		return [32]byte{}, err
	}
	return types.ComputeETHSigningRoot(r.ProposerPreferences, domain)
}

// carryOverPreferenceShares copies the prior sub-runner's collected preference shares for this sub's frozen
// preference root, so a re-emission that re-derives a byte-identical preference (same signing root) keeps its
// progress, as in carryOverAuthShares. When the preference changed, the roots differ and nothing carries.
func (r *ProposerPreferencesSlotRunner) carryOverPreferenceShares(prev *ProposerPreferencesSlotRunner) {
	if prev == nil || prev.GetState() == nil || r.GetState() == nil || r.ProposerPreferences == nil {
		return
	}
	signingRoot, err := r.preferenceSigningRoot()
	if err != nil {
		return
	}
	r.carryOverSharesForRoot(prev, signingRoot)
}

// carryOverSharesForRoot copies the prior sub-runner's collected partial signatures for one signing root
// into this sub's container. Callers guarantee both States are non-nil.
func (r *ProposerPreferencesSlotRunner) carryOverSharesForRoot(prev *ProposerPreferencesSlotRunner, signingRoot [32]byte) {
	validatorIndex := r.GetShare().ValidatorIndex
	for signer, sig := range prev.GetState().PreConsensusContainer.GetSignatures(validatorIndex, signingRoot) {
		r.GetState().PreConsensusContainer.AddSignature(&types.PartialSignatureMessage{
			PartialSignature: sig,
			SigningRoot:      signingRoot,
			Signer:           signer,
			ValidatorIndex:   validatorIndex,
		})
	}
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
