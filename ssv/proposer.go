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

	// producedEnvelope is this operator's own blinded execution-payload envelope from its Gloas
	// produceBlockV4 response (SIP #94 §6). It is what lets the operator publish the §6 reveal when it
	// turns out to be the builder operator — its BeaconBlockRoot equals the decided block root. Nil
	// pre-Gloas, on an external bid, until the Gloas produce in ProcessPreConsensus, and reset per duty
	// in executeDuty.
	//
	// Deliberately unexported: being the builder operator is a local, non-consensus property, so it stays
	// out of the JSON post-state root the spec vectors compare — unlike the other runners' per-duty state
	// (PayloadAttestationData, ProposerPreferences, BuilderRequestAuths), which is agreed and exported. A
	// builder and a non-builder therefore finish a §6 self-build duty with identical post-state; the
	// publish/don't-publish decision shows only in the beacon broadcast, covered by the runner tests
	// (full-flow builder publish, seeded non-builder no-publish) rather than a state-comparison vector.
	producedEnvelope *gloas.BlindedExecutionPayloadEnvelope

	// blockSubmitAttempted records that the decided block's signature reconstructed and its submit ran this
	// duty (regardless of the submit's result). The §6 reveal is gated on this rather than on a successful
	// submit, so it still publishes when this operator's own submit fails but the block reaches beacon nodes
	// from other operators' submits (SIP #94 §6 "retry until they have"; a node ignores an envelope whose
	// block it has not seen). Transient and unexported, reset per duty in executeDuty (after the duty is
	// accepted), like producedEnvelope.
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

	// get block data
	var input *types.ProposerConsensusData
	if versionForSlot(r.beacon, duty.Slot) >= gloas.DataVersionGloas {
		// Gloas (ePBS §4): a self-build produce yields the block and its own blinded envelope (SIP #94 §6),
		// so the operator carries payload_root in the decided value and — if it turns out to be the builder
		// operator — publishes the reveal. An external bid win (p2p/builder-API) returns a bare block and a
		// nil envelope, so payload_root is zero (the value check pins it zero iff not self-build).
		// api.VersionedProposal cannot carry a Gloas block, so the {block, payload_root} wrapper travels
		// opaque in DataSSZ (decoded by the value check and the post-consensus paths).
		blk, envelope, err := r.GetBeaconNode().GetGloasBeaconBlock(duty.Slot, r.GetShare().Graffiti, fullSig)
		if err != nil {
			return errors.Wrap(err, "failed to get Gloas Beacon block")
		}
		r.producedEnvelope = envelope
		var payloadRoot phase0.Root
		if envelope != nil {
			payloadRoot = envelope.PayloadRoot
		}
		byts, err := (&gloas.GloasProposalData{Block: blk, PayloadRoot: payloadRoot}).MarshalSSZ()
		if err != nil {
			return errors.Wrap(err, "could not marshal Gloas proposal data")
		}
		input = &types.ProposerConsensusData{
			Duty: *duty,
			// Stamp the slot's fork, not the Gloas constant, so the value check's cd.Version == slotVersion
			// pin holds at any fork the Gloas branch serves (SIP #94 §4); identical to the constant today.
			Version: versionForSlot(r.beacon, duty.Slot),
			DataSSZ: byts,
		}
	} else {
		vBlk, obj, err := r.GetBeaconNode().GetBeaconBlock(duty.Slot, r.GetShare().Graffiti, fullSig)
		if err != nil {
			return errors.Wrap(err, "failed to get Beacon block")
		}
		byts, err := obj.MarshalSSZ()
		if err != nil {
			return errors.Wrap(err, "could not marshal beacon block")
		}
		input = &types.ProposerConsensusData{
			Duty:    *duty,
			Version: vBlk.Version,
			DataSSZ: byts,
		}
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

	// Post-consensus entries: the block root under DomainProposer always, and — on the Gloas self-build
	// path — the §6 blinded-envelope root under DomainBeaconBuilder, riding the same packet (SIP #94 §4).
	var entries []*types.PartialSignatureMessage
	if versionForSlot(r.beacon, duty.Slot) >= gloas.DataVersionGloas {
		proposalData, err := gloas.DecodeGloasProposalData(cd.DataSSZ)
		if err != nil {
			return errors.Wrap(err, "could not decode Gloas proposal data from consensus data")
		}
		blockMsg, err := r.BaseRunner.signBeaconObject(r, duty, proposalData.Block, duty.Slot, types.DomainProposer)
		if err != nil {
			return errors.Wrap(err, "could not sign Gloas block")
		}
		entries = append(entries, blockMsg)

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
	} else {
		_, blkToSign, err := cd.GetBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get block data")
		}
		blockMsg, err := r.BaseRunner.signBeaconObject(r, duty, blkToSign, duty.Slot, types.DomainProposer)
		if err != nil {
			return errors.Wrap(err, "failed signing block")
		}
		entries = append(entries, blockMsg)
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
	// Fork and epoch come from the running duty's slot, not the decided value's own Duty.Slot (the two
	// agree once ProcessConsensus binds them); cd is used only for the block and envelope content.
	dutySlot := r.BaseRunner.State.StartingDuty.DutySlot()
	isGloas := versionForSlot(r.beacon, dutySlot) >= gloas.DataVersionGloas

	expectedRoots, err := r.expectedPostConsensusRootsAndDomains()
	if err != nil {
		return err
	}
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(dutySlot)
	blockSigningRoot, envelopeSigningRoot, err := r.classifyPostConsensusSigningRoots(expectedRoots, epoch)
	if err != nil {
		return err
	}

	// A Gloas self-build packet carries two roots (block + optional §6 envelope) that reconstruct
	// independently, so each is handled on its own: a failure on one must not abort the other, and a bad
	// share dropped by the fallback re-crosses in a later packet. The block (required) is handled first
	// regardless of packet order — it finalizes the duty, and a BN ignores an envelope whose block it has not
	// seen (SIP #94 §4/§6). When both fail the block error takes precedence (the block is the primary output;
	// the reveal is best-effort), so returning it over the envelope error loses nothing important.
	var blockErr, envelopeErr error
	if slices.Contains(roots, blockSigningRoot) {
		blockErr = r.submitDecidedBlock(cd, blockSigningRoot, isGloas)
	}
	// The optional §6 reveal publishes once the block submit has been attempted (blockSubmitAttempted — not
	// gated on submit success, so a failed local submit does not drop the reveal; SIP #94 §6) and the envelope
	// is at quorum. Gated on HasQuorum (not a first crossing) so an envelope quorum reached before the block is
	// not lost; awaitingEnvelope() stops admitting packets once it is at quorum, so this publishes exactly once.
	if isGloas && envelopeSigningRoot != ([32]byte{}) && r.blockSubmitAttempted &&
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
func (r *ProposerRunner) submitDecidedBlock(cd *types.ProposerConsensusData, root [32]byte, isGloas bool) error {
	sig, err := r.reconstructPostConsensusSig(root)
	if err != nil {
		return err
	}
	// The block signature reconstructed, so the submit below is an attempt regardless of its result — gate
	// the §6 reveal on this, not on submit success (see blockSubmitAttempted, SIP #94 §6).
	r.blockSubmitAttempted = true
	if isGloas {
		proposalData, err := gloas.DecodeGloasProposalData(cd.DataSSZ)
		if err != nil {
			return errors.Wrap(err, "could not decode Gloas proposal data from consensus data")
		}
		if err := r.GetBeaconNode().SubmitGloasBeaconBlock(proposalData.Block, sig); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed Gloas block")
		}
	} else {
		vBlk, _, err := cd.GetBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get block")
		}
		if err := r.GetBeaconNode().SubmitBeaconBlock(vBlk, sig); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed Beacon block")
		}
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

// classifyPostConsensusSigningRoots computes the expected block and §6 envelope signing roots (each under
// its domain) so ProcessPostConsensus can tell which quorum-reached root is which. envelope is zero when
// the packet carries no envelope entry (SIP #94 §4/§6).
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

// publishEnvelope publishes the §6 reveal on envelope-root quorum, but only for the builder operator — the
// one whose own produceBlockV4 response holds the decided block, i.e. its produced envelope equals the one
// derived from the decided value. Every other operator reconstructs the signature but publishes nothing
// (SIP #94 §6).
func (r *ProposerRunner) publishEnvelope(cd *types.ProposerConsensusData, sig phase0.BLSSignature) error {
	proposalData, err := gloas.DecodeGloasProposalData(cd.DataSSZ)
	if err != nil {
		return errors.Wrap(err, "could not decode Gloas proposal data from consensus data")
	}
	derived, err := deriveBlindedEnvelope(proposalData)
	if err != nil {
		return errors.Wrap(err, "could not derive blinded envelope")
	}
	if r.producedEnvelope == nil {
		return nil
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

// awaitingEnvelope reports whether the block is decided (so the duty is finished) but the optional §6
// envelope root is expected and has not yet reached post-consensus quorum — the window in which the
// proposer keeps accepting post-consensus packets (ValidatePostConsensusMsg) so an envelope quorum reached
// only after the block's still publishes the reveal (SIP #94 §4). It is derived from state on each call, so
// it is independent of packet order and needs no reset between duties: false pre-Gloas and on an external
// bid (no envelope root), with no decided value (a never-started duty), and once the envelope reconstructs.
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

// deriveBlindedEnvelope builds the §6 blinded envelope entirely from the decided GloasProposalData: with
// payload_root carried in the value, every field is derivable from the decided block (SIP #94 §6).
func deriveBlindedEnvelope(d *gloas.GloasProposalData) (*gloas.BlindedExecutionPayloadEnvelope, error) {
	blockRoot, err := d.Block.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not hash Gloas block")
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
// sign, each with its domain. At Gloas self-build slots it is the block root under DomainProposer plus the
// optional §6 envelope root under DomainBeaconBuilder (SIP #94 §4/§6).
func (r *ProposerRunner) expectedPostConsensusRootsAndDomains() ([]PostConsensusRoot, error) {
	cd := &types.ProposerConsensusData{}
	if err := cd.Decode(r.GetState().DecidedValue); err != nil {
		return nil, errors.Wrap(err, "could not create consensus data")
	}

	// Fork comes from the running duty's slot, not the decided value's own Duty.Slot; cd supplies only the
	// block and envelope content the roots are derived from.
	if versionForSlot(r.beacon, r.BaseRunner.State.StartingDuty.DutySlot()) >= gloas.DataVersionGloas {
		proposalData, err := gloas.DecodeGloasProposalData(cd.DataSSZ)
		if err != nil {
			return nil, errors.Wrap(err, "could not decode Gloas proposal data from consensus data")
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

	_, data, err := cd.GetBlockData()
	if err != nil {
		return nil, errors.Wrap(err, "could not get block data")
	}
	return SingleDomainPostConsensusRoots(types.DomainProposer, data), nil
}

// executeDuty steps:
// 1) sign a partial randao sig and wait for 2f+1 partial sigs from peers
// 2) reconstruct randao and send GetBeaconBlock to BN
// 3) start consensus on duty + block data
// 4) Once consensus decides, sign partial block and broadcast
// 5) collect 2f+1 partial sigs, reconstruct and broadcast valid block sig to the BN
func (r *ProposerRunner) executeDuty(duty types.Duty) error {
	// Reset the transient per-duty §6 state here, not in StartNewDuty: executeDuty runs only after the duty
	// is accepted, so a duty that ShouldProcessDuty rejects (a duplicate or a past slot) cannot wipe the
	// reveal state of the duty still in flight (SIP #94 §6).
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
