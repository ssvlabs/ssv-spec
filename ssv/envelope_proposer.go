package ssv

import (
	"bytes"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// EnvelopeProposerRunner runs the Gloas (ePBS) execution-payload envelope duty for the self-build
// proposer (SIP #94 §6). It has NO consensus phase: once §4 decides, bid.block_hash pins exactly one
// valid envelope, so there is nothing to negotiate. The flow is one dissemination round plus one
// threshold-signing round:
//
//  1. The builder operator — the one whose beacon node built the §4-decided block — disseminates the
//     blinded envelope of that block (SSVEnvelopeDisseminationMsgType).
//  2. Every operator content-selects the first disseminated envelope that binds to its own §4 decision
//     (the four checks in bindsToProposal), then threshold-signs the selected envelope's root under
//     DomainBeaconBuilder. The single signing round reuses the pre-consensus container, the same shape
//     as the PTC and proposer-preferences runners.
//  3. On quorum every operator reconstructs the signature; only the builder operator, which holds the
//     produced envelope matching the selected one, publishes the reveal.
type EnvelopeProposerRunner struct {
	BaseRunner *BaseRunner

	// ProposedBlocks is the §4→§6 linkage store (shared with the proposer runner in production): the
	// duty binds a disseminated envelope to the slot's §4-decided block recorded here.
	ProposedBlocks ProposedBlocks
	// ProducedEnvelope is this operator's own produced blinded envelope, held only by the builder
	// operator (nil otherwise). The publish-by-content-match compares it against the selected envelope.
	ProducedEnvelope *gloas.BlindedExecutionPayloadEnvelope
	// SelectedEnvelope is the disseminated envelope this operator chose to sign — the first arrival that
	// binds to its §4 decision. Incoming partial signatures validate against its root; nil until
	// selection, so peers' partials are rejected until then.
	SelectedEnvelope *gloas.BlindedExecutionPayloadEnvelope

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
}

func NewEnvelopeProposerRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	proposedBlocks ProposedBlocks,
) (Runner, error) {

	if len(share) != 1 {
		return nil, fmt.Errorf("must have one share")
	}
	if err := validateShareMap(share); err != nil {
		return nil, err
	}

	return &EnvelopeProposerRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RoleEnvelopeProposer,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},

		ProposedBlocks: proposedBlocks,

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
	}, nil
}

func (r *EnvelopeProposerRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	// Clear any prior envelope; executeDuty re-derives them, so a not-yet-executed or non-builder duty
	// stays nil.
	r.ProducedEnvelope = nil
	r.SelectedEnvelope = nil
	return r.BaseRunner.baseStartNewNonBeaconDuty(r, duty.(*types.ValidatorDuty), quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *EnvelopeProposerRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *EnvelopeProposerRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing envelope partial signature message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// Defensive: peer partials are validated against the selected envelope, so a quorum without one
	// should be unreachable.
	if r.SelectedEnvelope == nil {
		return types.NewError(types.EnvelopeNoSelectedEnvelopeErrorCode, "reached quorum without a selected envelope")
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

	// Publish by content match: only the builder operator holds the produced envelope that blinds to the
	// selected value, so only it can form a valid publish body; everyone else finishes without
	// publishing (SIP #94 §6).
	if r.builtSelectedEnvelope() {
		if err := r.beacon.SubmitExecutionPayloadEnvelope(r.SelectedEnvelope, specSig); err != nil {
			return errors.Wrap(err, "could not submit execution payload envelope")
		}
	}

	r.GetState().Finished = true
	return nil
}

func (r *EnvelopeProposerRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return types.NewError(types.EnvelopeProposerNoConsensusPhaseErrorCode, "no consensus phase for envelope proposer")
}

func (r *EnvelopeProposerRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.EnvelopeProposerNoPostConsensusPhaseErrorCode, "no post consensus phase for envelope proposer")
}

// ProcessEnvelopeDissemination handles a disseminated blinded envelope (SIP #94 §6): it selects the
// first arrival that binds to this operator's §4 decision and threshold-signs it. Non-binding
// disseminations are skipped silently (the binding checks are runner concerns, not validation rules, so
// they carry no peer penalty), and further disseminations are ignored once an envelope is selected.
func (r *EnvelopeProposerRunner) ProcessEnvelopeDissemination(msg *types.SignedSSVMessage) error {
	// Already selected and signed — ignore later arrivals (per-signer dedup is a validation-layer rule).
	if r.SelectedEnvelope != nil {
		return nil
	}

	dissemination := &types.EnvelopeDissemination{}
	if err := dissemination.Decode(msg.SSVMessage.Data); err != nil {
		return types.WrapError(types.EnvelopeDisseminationDecodeErrorCode, errors.Wrap(err, "could not decode envelope dissemination"))
	}

	duty, ok := r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty)
	if !ok {
		return types.NewError(types.InvalidValidatorDutyErrorCode, "starting duty is not a validator duty")
	}
	slot := duty.Slot

	// A dissemination for another slot is not this duty's; ignore it.
	if dissemination.Slot != slot {
		return nil
	}

	proposal, ok := r.ProposedBlocks.Get(slot)
	if !ok {
		// This operator has not decided §4 for the slot yet; it cannot validate the envelope. It keeps
		// running so a later arrival (after its block instance decides) can still be selected.
		return nil
	}

	// Content-based selection: sign the first disseminated envelope that binds to the §4 decision,
	// skipping any that fail (SIP #94 §6).
	if !bindsToProposal(dissemination.Envelope, proposal) {
		return nil
	}

	return r.selectAndSign(duty, slot, dissemination.Envelope)
}

func (r *EnvelopeProposerRunner) expectedPreConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	// Peer partials validate against this operator's selected envelope; with none selected — the operator
	// has not seen a binding dissemination — every peer message is rejected.
	if r.SelectedEnvelope == nil {
		return nil, types.DomainError, types.NewError(types.EnvelopeNoSelectedEnvelopeErrorCode, "no selected envelope")
	}
	return []types.HashRoot{r.SelectedEnvelope}, types.DomainBeaconBuilder, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *EnvelopeProposerRunner) expectedPostConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, fmt.Errorf("no post consensus roots for envelope proposer")
}

// executeDuty steps (SIP #94 §6):
//  1. read the slot's §4 decision — the duty binds against it; without it the operator waits
//  2. ask the local beacon node for this operator's own produced envelope; it answers only where the
//     block was built, so success identifies the builder operator
//  3. the builder operator disseminates its blinded envelope and, since its own envelope binds by
//     construction, selects and signs it; every other operator disseminates nothing and only signs a
//     dissemination it later receives (ProcessEnvelopeDissemination)
func (r *EnvelopeProposerRunner) executeDuty(duty types.Duty) error {
	r.ProducedEnvelope = nil
	r.SelectedEnvelope = nil
	slot := duty.DutySlot()

	proposal, ok := r.ProposedBlocks.Get(slot)
	if !ok {
		// No §4 decision recorded yet — the duty stays running and acts once a decision and a binding
		// dissemination arrive.
		return nil
	}

	envelope, err := r.beacon.GetBlindedExecutionPayloadEnvelope(slot, proposal.BlockRoot)
	if err != nil {
		// Not the builder operator (this beacon node did not build the decided block): disseminate
		// nothing and wait for a peer's dissemination.
		return nil
	}
	r.ProducedEnvelope = envelope

	if err := r.disseminateEnvelope(slot, envelope); err != nil {
		return errors.Wrap(err, "could not disseminate envelope")
	}

	// The builder operator's own envelope binds by construction; select and sign it now.
	if !bindsToProposal(envelope, proposal) {
		return nil
	}
	return r.selectAndSign(duty.(*types.ValidatorDuty), slot, envelope)
}

// selectAndSign records the selected envelope, signs its root under DomainBeaconBuilder, and broadcasts
// the partial signature for the threshold-signing round.
func (r *EnvelopeProposerRunner) selectAndSign(duty *types.ValidatorDuty, slot phase0.Slot, envelope *gloas.BlindedExecutionPayloadEnvelope) error {
	r.SelectedEnvelope = envelope

	msg, err := r.BaseRunner.signBeaconObject(r, duty, envelope, slot, types.DomainBeaconBuilder)
	if err != nil {
		return errors.Wrap(err, "could not sign blinded envelope")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.EnvelopePartialSig,
		Slot:     slot,
		Messages: []*types.PartialSignatureMessage{msg},
	}
	return r.broadcastPartialSig(msgs)
}

// disseminateEnvelope broadcasts the blinded envelope as an SSVEnvelopeDisseminationMsgType message.
func (r *EnvelopeProposerRunner) disseminateEnvelope(slot phase0.Slot, envelope *gloas.BlindedExecutionPayloadEnvelope) error {
	dissemination := &types.EnvelopeDissemination{
		Slot:     slot,
		Envelope: envelope,
	}
	data, err := dissemination.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode envelope dissemination")
	}

	msgID := types.NewValidatorMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.RunnerRoleType)
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVEnvelopeDisseminationMsgType,
		MsgID:   msgID,
		Data:    data,
	}
	return r.signAndBroadcast(ssvMsg)
}

// broadcastPartialSig wraps the partial-signature container in an SSVMessage and broadcasts it.
func (r *EnvelopeProposerRunner) broadcastPartialSig(msgs *types.PartialSignatureMessages) error {
	data, err := msgs.Encode()
	if err != nil {
		return err
	}
	msgID := types.NewValidatorMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.RunnerRoleType)
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    data,
	}
	return r.signAndBroadcast(ssvMsg)
}

// signAndBroadcast operator-signs an SSVMessage and broadcasts it.
func (r *EnvelopeProposerRunner) signAndBroadcast(ssvMsg *types.SSVMessage) error {
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
		return errors.Wrap(err, "can't broadcast envelope message")
	}
	return nil
}

// bindsToProposal reports whether a disseminated envelope matches the §4 decision (SIP #94 §6): the
// four checks that need no payload bytes. PayloadRoot is trusted from the builder operator (no local
// check), matching the blinded-block trust model.
func bindsToProposal(envelope *gloas.BlindedExecutionPayloadEnvelope, proposal ProposedBlock) bool {
	if envelope == nil || envelope.ExecutionRequests == nil {
		return false
	}
	if envelope.BuilderIndex != gloas.BuilderIndexSelfBuild {
		return false
	}
	if envelope.BeaconBlockRoot != proposal.BlockRoot {
		return false
	}
	if envelope.ParentBeaconBlockRoot != proposal.ParentRoot {
		return false
	}
	requestsRoot, err := envelope.ExecutionRequests.HashTreeRoot()
	if err != nil {
		return false
	}
	return requestsRoot == proposal.ExecutionRequestsRoot
}

// builtSelectedEnvelope reports whether this operator's own produced envelope is the selected one — the
// publish-by-content-match: only the builder operator holds the body behind the reconstructed signature.
func (r *EnvelopeProposerRunner) builtSelectedEnvelope() bool {
	if r.ProducedEnvelope == nil || r.SelectedEnvelope == nil {
		return false
	}
	produced, err := r.ProducedEnvelope.Encode()
	if err != nil {
		return false
	}
	selected, err := r.SelectedEnvelope.Encode()
	if err != nil {
		return false
	}
	return bytes.Equal(produced, selected)
}

func (r *EnvelopeProposerRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *EnvelopeProposerRunner) GetNetwork() Network {
	return r.network
}

func (r *EnvelopeProposerRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *EnvelopeProposerRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *EnvelopeProposerRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *EnvelopeProposerRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return nil
}

func (r *EnvelopeProposerRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *EnvelopeProposerRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}
