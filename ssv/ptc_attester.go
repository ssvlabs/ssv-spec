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

// PTCAttesterRunner runs the Gloas (ePBS) Payload Timeliness Committee attestation duty (SIP #94 §3).
// It has no consensus or post-consensus phase: each operator signs the PayloadAttestationData its own
// beacon node reports at execution time, and the per-validator signature reconstructs only once a
// quorum of operators converged on byte-identical data — honest convergence, not consensus.
type PTCAttesterRunner struct {
	BaseRunner *BaseRunner

	// PayloadAttestationData is the operator's frozen observation. Incoming partial signatures are
	// validated and aggregated against exactly this root; nil means the operator abstained (saw no
	// block for the slot) and is sitting the duty out.
	PayloadAttestationData *gloas.PayloadAttestationData

	// PendingPartials keeps, per slot, the partials received before that slot's duty started here, which is
	// when this operator's observation is frozen. StartNewDuty replays them (SIP #94 §3): a peer signs as soon
	// as its own observation is final, which can be earlier than this operator's, and its partial is delivered
	// once (§7 dedups repeats). A signer has one partial per slot (§7), so its first is the one kept, and a
	// slot's entry goes when a duty for that slot or a later one starts. The slot is the sender's choice and
	// can be ahead of the current one: what bounds the map is the §7 message-validation window on the slot,
	// which is outside the spec. Omitted from JSON when empty, so runners holding none keep their state roots.
	PendingPartials map[phase0.Slot][]*types.PartialSignatureMessages `json:"PendingPartials,omitempty"`

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
}

func NewPTCAttesterRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
) (Runner, error) {

	if len(share) != 1 {
		return nil, fmt.Errorf("must have one share")
	}
	if err := validateShareMap(share); err != nil {
		return nil, err
	}

	return &PTCAttesterRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RolePTCAttester,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},
		PendingPartials: map[phase0.Slot][]*types.PartialSignatureMessages{},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
	}, nil
}

// StartNewDuty starts the slot's duty, then replays the partials retained for it (SIP #94 §3).
func (r *PTCAttesterRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	// A rejected duty (a duplicate, or one for a past slot) leaves the retained partials as they are.
	if err := r.BaseRunner.ShouldProcessNonBeaconDuty(duty); err != nil {
		return errors.Wrap(err, "can't start non-beacon duty")
	}
	r.BaseRunner.baseSetupForNewDuty(duty, quorum)
	startErr := r.executeDuty(duty)

	// Replay each retained partial as if it just arrived. One that fails its checks is skipped: a diverging
	// peer's, or every one when this operator abstained or failed to observe. A failure to reconstruct or
	// submit a quorum it completes is returned, as on the live path, after the duty's own error.
	var replayErr error
	for _, msg := range r.takePendingPartials(duty.DutySlot()) {
		reachedQuorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, msg)
		if err != nil || !reachedQuorum {
			continue
		}
		// only 1 root, verified in basePreConsensusMsgProcessing
		if err := r.submitPayloadAttestation(roots[0]); err != nil && replayErr == nil {
			replayErr = err
		}
	}
	if startErr != nil {
		return startErr
	}
	return replayErr
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *PTCAttesterRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

// ProcessPreConsensus retains a partial whose slot's duty has not started here yet (PendingPartials), and
// collects any other against the running duty's frozen observation.
func (r *PTCAttesterRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	if r.awaitsDuty(signedMsg.Slot) {
		// Reject a malformed packet before retaining it, so it takes no room.
		if err := r.validatePartialsPacket(signedMsg); err != nil {
			return err
		}
		r.retainPartial(signedMsg)
		return nil
	}

	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing payload attestation message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// only 1 root, verified in basePreConsensusMsgProcessing
	return r.submitPayloadAttestation(roots[0])
}

// awaitsDuty reports whether no duty has started here yet for the slot or a later one, so the slot's
// observation is not frozen and its partials cannot be checked yet.
func (r *PTCAttesterRunner) awaitsDuty(slot phase0.Slot) bool {
	return r.BaseRunner.State == nil || r.BaseRunner.State.StartingDuty.DutySlot() < slot
}

// validatePartialsPacket applies the checks that need no frozen observation: the partials are valid, from one
// signer and for this runner's validator, and the packet carries the one partial of a PTC duty (SIP #94 §7).
// The running duty applies them too, and checks the root against its observation.
func (r *PTCAttesterRunner) validatePartialsPacket(signedMsg *types.PartialSignatureMessages) error {
	if err := signedMsg.Validate(); err != nil {
		return errors.Wrap(err, "PartialSignatureMessages invalid")
	}
	if err := r.BaseRunner.validateValidatorIndexInPartialSigMsg(signedMsg); err != nil {
		return err
	}
	if len(signedMsg.Messages) != 1 {
		return types.NewError(types.WrongRootsCountErrorCode, "ptc attester container must carry one partial")
	}
	return nil
}

// retainPartial adds the message to its slot's PendingPartials, unless its signer already has one there.
func (r *PTCAttesterRunner) retainPartial(signedMsg *types.PartialSignatureMessages) {
	// The map is omitted from JSON when empty, so a runner rebuilt from JSON holds a nil one.
	if r.PendingPartials == nil {
		r.PendingPartials = map[phase0.Slot][]*types.PartialSignatureMessages{}
	}
	signer := signedMsg.Messages[0].Signer
	for _, existing := range r.PendingPartials[signedMsg.Slot] {
		if existing.Messages[0].Signer == signer {
			return
		}
	}
	r.PendingPartials[signedMsg.Slot] = append(r.PendingPartials[signedMsg.Slot], signedMsg)
}

// takePendingPartials returns the partials retained for the slot and drops them along with those of earlier
// slots, whose duties can no longer start (ShouldProcessNonBeaconDuty).
func (r *PTCAttesterRunner) takePendingPartials(slot phase0.Slot) []*types.PartialSignatureMessages {
	pending := r.PendingPartials[slot]
	for s := range r.PendingPartials {
		if s <= slot {
			delete(r.PendingPartials, s)
		}
	}
	return pending
}

// submitPayloadAttestation reconstructs the signature over the frozen observation once its root reached
// quorum, and submits the PayloadAttestationMessage (SIP #94 §3).
func (r *PTCAttesterRunner) submitPayloadAttestation(root [32]byte) error {
	// Defensive: peer partials are validated against the frozen observation, so a quorum without one
	// should be unreachable.
	if r.PayloadAttestationData == nil {
		return types.NewError(types.PTCAttesterNoObservationErrorCode, "reached quorum without a frozen payload attestation data")
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

	msg := &gloas.PayloadAttestationMessage{
		ValidatorIndex: r.GetShare().ValidatorIndex,
		Data:           r.PayloadAttestationData,
		Signature:      specSig,
	}

	if err := r.beacon.SubmitPayloadAttestation(msg); err != nil {
		return errors.Wrap(err, "could not submit payload attestation message")
	}

	r.GetState().Finished = true
	return nil
}

func (r *PTCAttesterRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return types.NewError(types.PTCAttesterNoConsensusPhaseErrorCode, "no consensus phase for ptc attestation")
}

func (r *PTCAttesterRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.PTCAttesterNoPostConsensusPhaseErrorCode, "no post consensus phase for ptc attestation")
}

func (r *PTCAttesterRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	// Peer partials validate against the operator's own frozen observation (honest convergence): a
	// diverging peer's root fails the expected-root check (WrongSigningRootErrorCode). With no observation
	// — the operator abstained (its beacon node saw no block for the slot) or failed to observe — there is
	// no reference to judge peers against, and it does not reconstruct from their partials (SIP #94 §3), so
	// every peer partial gets PTCAttesterNoObservation. That is a divergence of observation ("my beacon node
	// didn't see the block and theirs did"), not a peer fault, so it is kept distinct from a wrong root:
	// §7 message validation never reads the local observation, and neither rejects nor ignores such a partial.
	if r.PayloadAttestationData == nil {
		return nil, types.DomainError, types.NewError(types.PTCAttesterNoObservationErrorCode, "no frozen payload attestation data")
	}
	return []ssz.HashRoot{r.PayloadAttestationData}, types.DomainPTCAttester, nil
}

// expectedPostConsensusRootsAndDomains an INTERNAL function, returns the expected post-consensus roots to sign
func (r *PTCAttesterRunner) expectedPostConsensusRootsAndDomains() ([]PostConsensusRoot, error) {
	return nil, fmt.Errorf("no post consensus roots for ptc attestation")
}

// executeDuty steps:
//  1. observe the slot's payload attestation data from the local beacon node; a zero BeaconBlockRoot
//     means no block was seen — abstain, freezing, signing and broadcasting nothing (SIP #94 §3)
//  2. freeze the observation, sign it under DomainPTCAttester and broadcast the partial signature
//  3. once a quorum of operators converged on the same data, reconstruct and submit the message
func (r *PTCAttesterRunner) executeDuty(duty types.Duty) error {
	// Clear the prior duty's observation only now that ShouldProcessNonBeaconDuty accepted this duty: a
	// rejected duplicate or past duty must not erase a running duty's observation, and an abstaining duty
	// must not inherit the prior one.
	r.PayloadAttestationData = nil

	slot := duty.DutySlot()

	data, err := r.beacon.GetPayloadAttestationData(slot)
	if err != nil {
		return errors.Wrap(err, "failed to get payload attestation data")
	}
	if data.BeaconBlockRoot == (phase0.Root{}) {
		// Abstain: the duty stays running with nothing frozen, so incoming peer partials are rejected.
		return nil
	}
	// SIP #94 §3 pins the observation's slot to the duty's. A BN answering with a different slot would
	// otherwise be signed under the duty slot's domain and broadcast in a duty-slot packet over data for
	// another slot; reject the mismatch rather than relabel it (a correct BN never trips this).
	if data.Slot != slot {
		return types.NewError(types.PTCAttesterWrongSlotErrorCode, "payload attestation data slot does not match duty slot")
	}

	r.PayloadAttestationData = data

	msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.ValidatorDuty), data, slot, types.DomainPTCAttester)
	if err != nil {
		return errors.Wrap(err, "could not sign payload attestation data")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.PTCAttesterPartialSig,
		Slot:     slot,
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
		return errors.Wrap(err, "can't broadcast partial payload attestation sig")
	}
	return nil
}

func (r *PTCAttesterRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *PTCAttesterRunner) GetNetwork() Network {
	return r.network
}

func (r *PTCAttesterRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *PTCAttesterRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *PTCAttesterRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *PTCAttesterRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return nil
}

func (r *PTCAttesterRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *PTCAttesterRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}
