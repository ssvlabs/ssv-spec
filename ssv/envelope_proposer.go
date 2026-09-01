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
// proposer (SIP #94 §6): a second QBFT instance over the blinded envelope of the §4-decided block,
// followed by post-consensus threshold signing under DomainBeaconBuilder. Only the operator whose own
// produced envelope byte-matches the decided value publishes it (it holds the body); the others finish
// without publishing.
//
// The runner models the blinded (stateful) publication path only: the blinded envelope's root equals
// the full envelope's, so the signature is valid for either form, and the stateless full-body variant
// is node-side plumbing.
type EnvelopeProposerRunner struct {
	BaseRunner *BaseRunner

	// ProposedBlockRoots is the §4→§6 linkage store (shared with the proposer runner in production):
	// the duty only exists for a slot whose §4-decided block root was recorded, and the value check
	// pins the envelope to it.
	ProposedBlockRoots ProposedBlockRoots
	// ProducedEnvelope is this operator's own produced blinded envelope for the running duty; the
	// publish-by-content-match compares the decided value against exactly it. Nil until the duty
	// executes.
	ProducedEnvelope *gloas.BlindedExecutionPayloadEnvelope

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewEnvelopeProposerRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	proposedBlockRoots ProposedBlockRoots,
	highestDecidedSlot phase0.Slot,
) (Runner, error) {

	if len(share) != 1 {
		return nil, fmt.Errorf("must have one share")
	}
	if err := validateShareMap(share); err != nil {
		return nil, err
	}

	return &EnvelopeProposerRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType:     types.RoleEnvelopeProposer,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		ProposedBlockRoots: proposedBlockRoots,

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		valCheck:       valCheck,
	}, nil
}

func (r *EnvelopeProposerRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	return r.BaseRunner.baseStartNewDuty(r, duty, quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *EnvelopeProposerRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *EnvelopeProposerRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	return types.NewError(types.EnvelopeProposerNoPreConsensusPhaseErrorCode, "no pre consensus phase for envelope proposer")
}

func (r *EnvelopeProposerRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg, &types.EnvelopeConsensusData{})
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	cd := decidedValue.(*types.EnvelopeConsensusData)
	blinded := &gloas.BlindedExecutionPayloadEnvelope{}
	if err := blinded.Decode(cd.DataSSZ); err != nil {
		return errors.Wrap(err, "could not decode blinded envelope from consensus data")
	}

	msg, err := r.BaseRunner.signBeaconObject(r, r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty), blinded,
		cd.Duty.Slot,
		types.DomainBeaconBuilder)
	if err != nil {
		return errors.Wrap(err, "failed signing blinded envelope")
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     cd.Duty.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
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

func (r *EnvelopeProposerRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
		if err != nil {
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			for _, root := range roots {
				r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PostConsensusContainer, root,
					r.GetShare().Committee, r.GetShare().ValidatorIndex)
			}
			return errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		cd := &types.EnvelopeConsensusData{}
		if err := cd.Decode(r.GetState().DecidedValue); err != nil {
			return errors.Wrap(err, "could not create consensus data")
		}

		// Publish by content match: only the operator whose own produced envelope is the decided one
		// holds the body to publish; everyone else finishes without publishing (SIP #94 §6).
		if r.builtDecidedEnvelope(cd.DataSSZ) {
			signed := &gloas.SignedBlindedExecutionPayloadEnvelope{
				Message:   r.ProducedEnvelope,
				Signature: specSig,
			}
			if err := r.GetBeaconNode().SubmitBlindedExecutionPayloadEnvelope(signed); err != nil {
				return errors.Wrap(err, "could not submit execution payload envelope")
			}
		}
	}
	r.GetState().Finished = true
	return nil
}

// builtDecidedEnvelope reports whether this operator's own produced envelope is the decided one
func (r *EnvelopeProposerRunner) builtDecidedEnvelope(decidedDataSSZ []byte) bool {
	if r.ProducedEnvelope == nil {
		return false
	}
	producedSSZ, err := r.ProducedEnvelope.Encode()
	if err != nil {
		return false
	}
	return bytes.Equal(producedSSZ, decidedDataSSZ)
}

func (r *EnvelopeProposerRunner) expectedPreConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	return nil, types.DomainError, types.NewError(types.EnvelopeProposerNoPreConsensusPhaseErrorCode, "no pre consensus phase for envelope proposer")
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *EnvelopeProposerRunner) expectedPostConsensusRootsAndDomain() ([]types.HashRoot, phase0.DomainType, error) {
	cd := &types.EnvelopeConsensusData{}
	if err := cd.Decode(r.GetState().DecidedValue); err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not create consensus data")
	}
	blinded := &gloas.BlindedExecutionPayloadEnvelope{}
	if err := blinded.Decode(cd.DataSSZ); err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not decode blinded envelope from consensus data")
	}
	return []types.HashRoot{blinded}, types.DomainBeaconBuilder, nil
}

// executeDuty steps:
//  1. read the slot's §4-decided block root — the duty only exists after the proposer decided (SIP #94 §6)
//  2. produce this operator's blinded envelope for that block and start consensus over it
//  3. once decided, sign the decided envelope under DomainBeaconBuilder (post-consensus), and on
//     quorum the operator that produced the decided envelope publishes it
func (r *EnvelopeProposerRunner) executeDuty(duty types.Duty) error {
	r.ProducedEnvelope = nil // drop any envelope produced for a prior duty
	slot := duty.DutySlot()

	blockRoot, ok := r.ProposedBlockRoots.Get(slot)
	if !ok {
		return types.NewError(types.EnvelopeNoProposedBlockRootErrorCode, "no decided block root recorded for the envelope slot")
	}

	envelope, err := r.GetBeaconNode().GetBlindedExecutionPayloadEnvelope(slot, blockRoot)
	if err != nil {
		return errors.Wrap(err, "failed to get blinded execution payload envelope")
	}
	r.ProducedEnvelope = envelope

	byts, err := envelope.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode blinded envelope")
	}
	input := &types.EnvelopeConsensusData{
		Duty:    *duty.(*types.ValidatorDuty),
		Version: gloas.DataVersionGloas,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, slot, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
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
	return r.valCheck
}

func (r *EnvelopeProposerRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *EnvelopeProposerRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}
