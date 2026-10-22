package ssv

import (
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

	// Post-consensus entries: the block root under DomainProposer always, and — on the self-build path — the
	// §6 blinded-envelope root under DomainBeaconBuilder, riding the same packet (SIP #94 §4).
	proposalData, err := cd.GetBlockData()
	if err != nil {
		return errors.Wrap(err, "could not get block data")
	}
	blockMsg, err := r.BaseRunner.signBeaconObject(r, duty, proposalData.Block, duty.Slot, types.DomainProposer)
	if err != nil {
		return errors.Wrap(err, "failed signing block")
	}
	entries := []*types.PartialSignatureMessage{blockMsg}

	if proposalData.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex == gloas.BuilderIndexSelfBuild {
		envelope, err := deriveBlindedEnvelope(proposalData)
		if err != nil {
			return errors.Wrap(err, "could not derive blinded envelope")
		}
		envMsg, err := r.BaseRunner.signBeaconObject(r, duty, envelope, duty.Slot, types.DomainBeaconBuilder)
		if err != nil {
			return errors.Wrap(err, "could not sign envelope")
		}
		entries = append(entries, envMsg)
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
	return nil
}

func (r *ProposerRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	cd := &types.ProposerConsensusData{}
	if err := cd.Decode(r.GetState().DecidedValue); err != nil {
		return errors.Wrap(err, "could not create consensus data")
	}
	expectedRoots, err := r.expectedPostConsensusRootsAndDomains()
	if err != nil {
		return err
	}
	// The epoch comes from the running duty's slot; cd supplies only the block and envelope content.
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.DutySlot())
	blockSigningRoot, envelopeSigningRoot, err := r.classifyPostConsensusSigningRoots(expectedRoots, epoch)
	if err != nil {
		return err
	}

	// The block and the optional §6 envelope reconstruct independently: a failure on one never aborts the other,
	// and a share the fallback drops re-crosses quorum in a later packet. The block goes first whatever the packet
	// order, since it finalizes the duty and a beacon node ignores an envelope whose block it has not seen
	// (SIP #94 §4/§6). If both fail, the block error is returned.
	var blockErr, envelopeErr error
	if slices.Contains(roots, blockSigningRoot) {
		blockErr = r.submitDecidedBlock(cd, blockSigningRoot)
	}
	// The reveal publishes once the block submit was attempted and the envelope is at quorum. HasQuorum (rather
	// than a first crossing) keeps an envelope quorum reached before the block; awaitingEnvelope stops admitting
	// packets once the envelope is at quorum, so the reveal publishes once.
	if envelopeSigningRoot != ([32]byte{}) && r.blockSubmitAttempted &&
		r.GetState().PostConsensusContainer.HasQuorum(r.GetShare().ValidatorIndex, envelopeSigningRoot) {
		envelopeErr = r.publishDecidedEnvelope(cd, envelopeSigningRoot)
	}
	if blockErr != nil {
		return blockErr
	}
	return envelopeErr
}

// reconstructPostConsensusSig reconstructs the threshold signature over root from the post-consensus
// container. On failure it falls back to verifying each partial and removing the invalid ones, so the
// root can re-cross quorum in a later packet once the offending share is gone.
func (r *ProposerRunner) reconstructPostConsensusSig(root [32]byte) (phase0.BLSSignature, error) {
	sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
	if err != nil {
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PostConsensusContainer, root,
			r.GetShare().Committee, r.GetShare().ValidatorIndex)
		return phase0.BLSSignature{}, errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
	}
	var specSig phase0.BLSSignature
	copy(specSig[:], sig)
	return specSig, nil
}

// submitDecidedBlock reconstructs the block signature and submits the decided block, finalizing the duty.
func (r *ProposerRunner) submitDecidedBlock(cd *types.ProposerConsensusData, root [32]byte) error {
	sig, err := r.reconstructPostConsensusSig(root)
	if err != nil {
		return err
	}
	// The §6 reveal gates on this attempt, not on the submit's result (see blockSubmitAttempted).
	r.blockSubmitAttempted = true
	proposalData, err := cd.GetBlockData()
	if err != nil {
		return errors.Wrap(err, "could not get block data")
	}
	if err := r.GetBeaconNode().SubmitBeaconBlock(proposalData.Block, sig); err != nil {
		return types.WrapError(types.ProposerBlockSubmitFailedErrorCode, errors.Wrap(err, "could not submit to Beacon chain reconstructed signed Beacon block"))
	}
	r.GetState().Finished = true
	return nil
}

// publishDecidedEnvelope reconstructs the §6 envelope signature and publishes the reveal (builder operator
// only; see publishEnvelope).
func (r *ProposerRunner) publishDecidedEnvelope(cd *types.ProposerConsensusData, root [32]byte) error {
	sig, err := r.reconstructPostConsensusSig(root)
	if err != nil {
		return err
	}
	return r.publishEnvelope(cd, sig)
}

// classifyPostConsensusSigningRoots computes the expected block and §6 envelope signing roots, each under its
// domain; envelope is zero when the decided value expects no envelope entry (an external bid).
func (r *ProposerRunner) classifyPostConsensusSigningRoots(expected []PostConsensusRoot, epoch phase0.Epoch) (block, envelope [32]byte, err error) {
	for _, e := range expected {
		d, err := r.GetBeaconNode().DomainData(epoch, e.Domain)
		if err != nil {
			return block, envelope, errors.Wrap(err, "could not get post-consensus root domain")
		}
		sr, err := types.ComputeETHSigningRoot(e.Root, d)
		if err != nil {
			return block, envelope, errors.Wrap(err, "could not compute ETH signing root")
		}
		switch e.Domain {
		case types.DomainProposer:
			block = sr
		case types.DomainBeaconBuilder:
			envelope = sr
		}
	}
	return block, envelope, nil
}

// publishEnvelope publishes the §6 reveal, but only for the builder operator: the one whose own produced
// envelope equals the one derived from the decided value. Every other operator publishes nothing (SIP #94 §6).
func (r *ProposerRunner) publishEnvelope(cd *types.ProposerConsensusData, sig phase0.BLSSignature) error {
	if r.producedEnvelope == nil {
		return nil
	}
	proposalData, err := cd.GetBlockData()
	if err != nil {
		return errors.Wrap(err, "could not get block data")
	}
	derived, err := deriveBlindedEnvelope(proposalData)
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
// (SIP #94 §4). Derived from state on each call, so it needs no per-duty reset.
func (r *ProposerRunner) awaitingEnvelope() bool {
	state := r.GetState()
	if state == nil || len(state.DecidedValue) == 0 {
		return false
	}
	expected, err := r.expectedPostConsensusRootsAndDomains()
	if err != nil {
		return false
	}
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(state.StartingDuty.DutySlot())
	for _, e := range expected {
		if !e.Optional {
			continue
		}
		d, err := r.GetBeaconNode().DomainData(epoch, e.Domain)
		if err != nil {
			return false
		}
		sr, err := types.ComputeETHSigningRoot(e.Root, d)
		if err != nil {
			return false
		}
		if !state.PostConsensusContainer.HasQuorum(r.GetShare().ValidatorIndex, sr) {
			return true
		}
	}
	return false
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

// expectedPostConsensusRootsAndDomains an INTERNAL function, returns the expected post-consensus roots to
// sign: the block root under DomainProposer and, on self-build, the optional §6 envelope root under
// DomainBeaconBuilder (SIP #94 §4/§6).
func (r *ProposerRunner) expectedPostConsensusRootsAndDomains() ([]PostConsensusRoot, error) {
	cd := &types.ProposerConsensusData{}
	if err := cd.Decode(r.GetState().DecidedValue); err != nil {
		return nil, errors.Wrap(err, "could not create consensus data")
	}
	proposalData, err := cd.GetBlockData()
	if err != nil {
		return nil, errors.Wrap(err, "could not get block data")
	}

	roots := []PostConsensusRoot{{Root: proposalData.Block, Domain: types.DomainProposer}}
	if proposalData.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex == gloas.BuilderIndexSelfBuild {
		envelope, err := deriveBlindedEnvelope(proposalData)
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
