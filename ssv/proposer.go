package ssv

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

type ProposerRunner struct {
	BaseRunner *BaseRunner

	// producedEnvelope is the blinded envelope from this operator's own self-build produce (SIP #94 §6). Its
	// BeaconBlockRoot matching the decided block makes this operator the builder operator, which publishes the
	// reveal. Nil on an external bid; reset per duty in executeDuty. Unexported on purpose: being the builder
	// operator is local, not agreed, state, so it stays out of the post-state root the vectors compare (the
	// publish decision shows in their beacon broadcasts).
	producedEnvelope *gloas.BlindedExecutionPayloadEnvelope

	// blockSubmitAttempted records that this duty's block signature reconstructed and its submit ran, whatever
	// the result. The §6 reveal gates on it rather than on a successful submit, since the block still reaches
	// beacon nodes through the other operators (SIP #94 §6). Reset per duty in executeDuty.
	blockSubmitAttempted bool

	// decided caches what the post-consensus round derives from the decided value (see decidedPostConsensus).
	// Unexported: derived state, kept out of the post-state root. It is keyed by duty slot and value, so it
	// needs no per-duty reset.
	decided *decidedProposal

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewProposerRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) (Runner, error) {

	if len(share) != 1 {
		return nil, fmt.Errorf("must have one share")
	}
	if err := validateShareMap(share); err != nil {
		return nil, err
	}

	return &ProposerRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType:     types.RoleProposer,
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
	}, nil
}

func (r *ProposerRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	return r.BaseRunner.baseStartNewDuty(r, duty, quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *ProposerRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *ProposerRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing randao message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// only 1 root, verified in basePreConsensusMsgProcessing
	root := roots[0]
	// randao is relevant only for block proposals, no need to check type
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
	if err != nil {
		// If the reconstructed signature verification failed, fall back to verifying each partial signature
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee,
			r.GetShare().ValidatorIndex)
		return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
	}

	duty := r.GetState().StartingDuty.(*types.ValidatorDuty)

	// A self-build produce also yields this operator's blinded envelope, whose payload_root rides in the decided
	// value (SIP #94 §4/§6); an external bid win returns a bare block and no envelope, so payload_root is zero.
	blk, envelope, err := r.GetBeaconNode().GetBeaconBlock(duty.Slot, r.GetShare().Graffiti, fullSig)
	if err != nil {
		return errors.Wrap(err, "failed to get Beacon block")
	}
	r.producedEnvelope = envelope
	var payloadRoot phase0.Root
	if envelope != nil {
		payloadRoot = envelope.PayloadRoot
	}
	byts, err := (&gloas.GloasProposalData{Block: blk, PayloadRoot: payloadRoot}).MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal proposal data")
	}
	input := &types.ProposerConsensusData{
		Duty:    *duty,
		Version: gloas.DataVersionGloas,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, input.Duty.DutySlot(), input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (r *ProposerRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg, &types.ProposerConsensusData{})
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	cd := decidedValue.(*types.ProposerConsensusData)
	duty := r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty)

	// Post-consensus entries, each under its own domain: the block root always and, on the self-build path,
	// the §6 blinded-envelope root, riding the same packet (SIP #94 §4).
	proposal, err := cd.GetBlockData()
	if err != nil {
		return errors.Wrap(err, "could not get block data")
	}
	roots, err := postConsensusRoots(proposal)
	if err != nil {
		return err
	}
	// A required root that fails to sign fails the round. An optional one (the §6 envelope root) is left out of
	// the packet instead: receivers do not require that entry, so the failure costs this operator its envelope
	// share and never its block share (SIP #94 §4). Its error is returned once the packet is broadcast.
	entries := make([]*types.PartialSignatureMessage, 0, len(roots))
	var envelopeSignErr error
	for _, root := range roots {
		msg, err := r.BaseRunner.signBeaconObject(r, duty, root.Root, duty.Slot, root.Domain)
		if err != nil {
			err = errors.Wrap(err, "failed signing post-consensus root")
			if !root.Optional {
				return err
			}
			envelopeSignErr = types.WrapError(types.ProposerEnvelopeSignFailedErrorCode, err)
			continue
		}
		entries = append(entries, msg)
	}

	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     duty.Slot,
		Messages: entries,
	}

	msgID := types.NewValidatorMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.RunnerRoleType)

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
	return envelopeSignErr
}

func (r *ProposerRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	decided, err := r.decidedPostConsensus()
	if err != nil {
		return err
	}

	// The block and the optional §6 envelope reconstruct independently: a failure on one never aborts the other,
	// and each recovers from a bad share on its own (see reconstructPostConsensusSig). The block goes first
	// whatever the packet order, since it finalizes the duty and a beacon node ignores an envelope whose block it
	// has not seen (SIP #94 §4/§6). If both fail, the block error is returned.
	var blockErr, envelopeErr error
	if slices.Contains(roots, decided.blockRoot) {
		blockErr = r.submitDecidedBlock(decided.proposal, decided.blockRoot)
	}
	// The reveal publishes once the block submit was attempted and the envelope is at quorum. HasQuorum (rather
	// than a first crossing) keeps an envelope quorum reached before the block; awaitingEnvelope stops admitting
	// packets once the envelope is at quorum, so the reveal publishes once.
	if decided.envelopeRoot != ([32]byte{}) && r.blockSubmitAttempted &&
		r.GetState().PostConsensusContainer.HasQuorum(r.GetShare().ValidatorIndex, decided.envelopeRoot) {
		envelopeErr = r.publishDecidedEnvelope(decided.proposal, decided.envelopeRoot)
	}
	if blockErr != nil {
		return blockErr
	}
	return envelopeErr
}

// reconstructPostConsensusSig reconstructs the threshold signature over root from the post-consensus
// container. On failure it verifies each partial and drops the invalid ones. If the valid ones are still a
// quorum it reconstructs from them, as no later packet would cross that quorum again; otherwise the root
// re-crosses quorum in a later packet.
func (r *ProposerRunner) reconstructPostConsensusSig(root [32]byte) (phase0.BLSSignature, error) {
	container := r.GetState().PostConsensusContainer
	validatorIndex := r.GetShare().ValidatorIndex
	sig, err := r.GetState().ReconstructBeaconSig(container, root, r.GetShare().ValidatorPubKey[:], validatorIndex)
	if err != nil {
		r.BaseRunner.FallBackAndVerifyEachSignature(container, root, r.GetShare().Committee, validatorIndex)
		if container.HasQuorum(validatorIndex, root) {
			sig, err = r.GetState().ReconstructBeaconSig(container, root, r.GetShare().ValidatorPubKey[:], validatorIndex)
		}
		if err != nil {
			return phase0.BLSSignature{}, errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
		}
	}
	var specSig phase0.BLSSignature
	copy(specSig[:], sig)
	return specSig, nil
}

// submitDecidedBlock reconstructs the block signature and submits the decided block, finalizing the duty.
func (r *ProposerRunner) submitDecidedBlock(proposal *gloas.GloasProposalData, root [32]byte) error {
	sig, err := r.reconstructPostConsensusSig(root)
	if err != nil {
		return err
	}
	// The §6 reveal gates on this attempt, not on the submit's result (see blockSubmitAttempted).
	r.blockSubmitAttempted = true
	if err := r.GetBeaconNode().SubmitBeaconBlock(proposal.Block, sig); err != nil {
		return types.WrapError(types.ProposerBlockSubmitFailedErrorCode, errors.Wrap(err, "could not submit to Beacon chain reconstructed signed Beacon block"))
	}
	r.GetState().Finished = true
	return nil
}

// publishDecidedEnvelope reconstructs the §6 envelope signature and publishes the reveal (builder operator
// only; see publishEnvelope).
func (r *ProposerRunner) publishDecidedEnvelope(proposal *gloas.GloasProposalData, root [32]byte) error {
	sig, err := r.reconstructPostConsensusSig(root)
	if err != nil {
		return err
	}
	return r.publishEnvelope(proposal, sig)
}

// decidedProposal is what the post-consensus round derives from a decided value: the proposal it carries, and
// the signing roots the round collects shares for.
type decidedProposal struct {
	// value and slot are the decided value and the duty slot the rest was derived for.
	value []byte
	slot  phase0.Slot

	proposal *gloas.GloasProposalData
	// blockRoot is the block's signing root under DomainProposer, and envelopeRoot the §6 envelope's under
	// DomainBeaconBuilder, both for the duty's epoch. envelopeRoot is zero when the proposal expects no envelope
	// entry (an external bid).
	blockRoot, envelopeRoot [32]byte
}

// decidedPostConsensus decodes the decided value and computes the signing roots the post-consensus round
// collects shares for. Packet validation, quorum handling and awaitingEnvelope all read its result, which is
// cached while the decided value and the duty are unchanged.
func (r *ProposerRunner) decidedPostConsensus() (*decidedProposal, error) {
	state := r.GetState()
	slot := state.StartingDuty.DutySlot()
	if r.decided != nil && r.decided.slot == slot && bytes.Equal(r.decided.value, state.DecidedValue) {
		return r.decided, nil
	}

	cd := &types.ProposerConsensusData{}
	if err := cd.Decode(state.DecidedValue); err != nil {
		return nil, errors.Wrap(err, "could not create consensus data")
	}
	proposal, err := cd.GetBlockData()
	if err != nil {
		return nil, errors.Wrap(err, "could not get block data")
	}
	expectedRoots, err := postConsensusRoots(proposal)
	if err != nil {
		return nil, err
	}

	decided := &decidedProposal{value: state.DecidedValue, slot: slot, proposal: proposal}
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(slot)
	for _, expectedRoot := range expectedRoots {
		domain, err := r.GetBeaconNode().DomainData(epoch, expectedRoot.Domain)
		if err != nil {
			return nil, errors.Wrap(err, "could not get post-consensus root domain")
		}
		signingRoot, err := types.ComputeETHSigningRoot(expectedRoot.Root, domain)
		if err != nil {
			return nil, errors.Wrap(err, "could not compute ETH signing root")
		}
		switch expectedRoot.Domain {
		case types.DomainProposer:
			decided.blockRoot = signingRoot
		case types.DomainBeaconBuilder:
			decided.envelopeRoot = signingRoot
		}
	}
	r.decided = decided
	return decided, nil
}

// publishEnvelope publishes the §6 reveal, but only for the builder operator: the one whose own produced
// envelope equals the one derived from the decided value. Every other operator publishes nothing (SIP #94 §6).
func (r *ProposerRunner) publishEnvelope(proposal *gloas.GloasProposalData, sig phase0.BLSSignature) error {
	if r.producedEnvelope == nil {
		return nil
	}
	derived, err := deriveBlindedEnvelope(proposal)
	if err != nil {
		return errors.Wrap(err, "could not derive blinded envelope")
	}
	produced, err := r.producedEnvelope.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not hash produced envelope")
	}
	want, err := derived.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not hash derived envelope")
	}
	if produced != want {
		// Not the builder operator: this operator's own produce is for a different block.
		return nil
	}
	if err := r.GetBeaconNode().SubmitExecutionPayloadEnvelope(r.producedEnvelope, sig); err != nil {
		return errors.Wrap(err, "could not submit execution payload envelope")
	}
	return nil
}

// awaitingEnvelope reports whether the decided value expects a §6 envelope root that has not reached
// post-consensus quorum yet. While it holds, a finished proposer keeps accepting post-consensus packets
// (ValidatePostConsensusMsg), so an envelope quorum completing after the block's still publishes the reveal
// (SIP #94 §4). It reads the state on each call, so it needs no per-duty reset.
func (r *ProposerRunner) awaitingEnvelope() bool {
	state := r.GetState()
	if state == nil || len(state.DecidedValue) == 0 {
		return false
	}
	decided, err := r.decidedPostConsensus()
	if err != nil || decided.envelopeRoot == ([32]byte{}) {
		return false
	}
	return !state.PostConsensusContainer.HasQuorum(r.GetShare().ValidatorIndex, decided.envelopeRoot)
}

// deriveBlindedEnvelope builds the §6 blinded envelope from the decided value: the block plus payload_root
// determine every field (SIP #94 §6).
func deriveBlindedEnvelope(d *gloas.GloasProposalData) (*gloas.BlindedExecutionPayloadEnvelope, error) {
	blockRoot, err := d.Block.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not hash block")
	}
	return &gloas.BlindedExecutionPayloadEnvelope{
		PayloadRoot:           d.PayloadRoot,
		ExecutionRequestsRoot: d.Block.Body.SignedExecutionPayloadBid.Message.ExecutionRequestsRoot,
		BuilderIndex:          gloas.BuilderIndexSelfBuild,
		BeaconBlockRoot:       blockRoot,
		ParentBeaconBlockRoot: d.Block.ParentRoot,
	}, nil
}

func (r *ProposerRunner) expectedPreConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.GetState().StartingDuty.DutySlot())
	return []types.HashRoot{types.SSZUint64(epoch)}, types.DomainRandao, nil
}

// expectedPostConsensusSigningRoots an INTERNAL function, returns the signing roots of the decided proposal's
// post-consensus roots: the block's, which a packet must carry, and on self-build the §6 envelope's, which it
// may (see postConsensusRoots).
func (r *ProposerRunner) expectedPostConsensusSigningRoots() ([]PostConsensusSigningRoot, error) {
	decided, err := r.decidedPostConsensus()
	if err != nil {
		return nil, err
	}
	expectedRoots := []PostConsensusSigningRoot{{Root: decided.blockRoot}}
	if decided.envelopeRoot != ([32]byte{}) {
		expectedRoots = append(expectedRoots, PostConsensusSigningRoot{Root: decided.envelopeRoot, Optional: true})
	}
	return expectedRoots, nil
}

// PostConsensusRoot is an object the proposer's post-consensus round signs, the domain to sign it under, and
// whether the round can do without it (see postConsensusRoots).
type PostConsensusRoot struct {
	Root     types.HashRoot
	Domain   phase0.DomainType
	Optional bool
}

// postConsensusRoots returns the roots the post-consensus round signs for a proposal: the block root under
// DomainProposer and, on self-build, the optional §6 envelope root under DomainBeaconBuilder
// (SIP #94 §4/§6).
func postConsensusRoots(proposal *gloas.GloasProposalData) ([]PostConsensusRoot, error) {
	roots := []PostConsensusRoot{{Root: proposal.Block, Domain: types.DomainProposer}}
	if proposal.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex == gloas.BuilderIndexSelfBuild {
		envelope, err := deriveBlindedEnvelope(proposal)
		if err != nil {
			return nil, errors.Wrap(err, "could not derive blinded envelope")
		}
		roots = append(roots, PostConsensusRoot{Root: envelope, Domain: types.DomainBeaconBuilder, Optional: true})
	}
	return roots, nil
}

// executeDuty steps:
// 1) sign a partial randao sig and wait for 2f+1 partial sigs from peers
// 2) reconstruct randao and send GetBeaconBlock to BN
// 3) start consensus on duty + block data
// 4) Once consensus decides, sign partial block and broadcast
// 5) collect 2f+1 partial sigs, reconstruct and broadcast valid block sig to the BN
func (r *ProposerRunner) executeDuty(duty types.Duty) error {
	// Reset the per-duty §6 state here rather than in StartNewDuty: executeDuty runs only once the duty is
	// accepted, so a rejected duplicate or past-slot duty cannot wipe the state of the duty in flight.
	r.blockSubmitAttempted = false
	r.producedEnvelope = nil

	// sign partial randao
	epoch := r.GetBeaconNode().GetBeaconNetwork().EstimatedEpochAtSlot(duty.DutySlot())
	msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.ValidatorDuty), types.SSZUint64(epoch), duty.DutySlot(),
		types.DomainRandao)
	if err != nil {
		return errors.Wrap(err, "could not sign randao")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     duty.DutySlot(),
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
		return errors.Wrap(err, "can't broadcast partial randao sig")
	}
	return nil
}

func (r *ProposerRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *ProposerRunner) GetNetwork() Network {
	return r.network
}

func (r *ProposerRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *ProposerRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *ProposerRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *ProposerRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *ProposerRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *ProposerRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}
