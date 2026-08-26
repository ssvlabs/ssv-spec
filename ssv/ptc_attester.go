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

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
	}, nil
}

func (r *PTCAttesterRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	// Clear any prior observation; executeDuty re-freezes it only if this operator attests, so an
	// abstained or not-yet-executed duty stays nil.
	r.PayloadAttestationData = nil
	return r.BaseRunner.baseStartNewNonBeaconDuty(r, duty.(*types.ValidatorDuty), quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *PTCAttesterRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *PTCAttesterRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing payload attestation message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// Defensive: peer partials are validated against the frozen observation, so a quorum without one
	// should be unreachable.
	if r.PayloadAttestationData == nil {
		return types.NewError(types.PTCAttesterNoObservationErrorCode, "reached quorum without a frozen payload attestation data")
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
	// diverging peer's root fails the expected-root check, and with no observation — the operator
	// abstained or has not executed the duty — every peer message is rejected.
	if r.PayloadAttestationData == nil {
		return nil, types.DomainError, types.NewError(types.PTCAttesterNoObservationErrorCode, "no frozen payload attestation data")
	}
	return []ssz.HashRoot{r.PayloadAttestationData}, types.DomainPTCAttester, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *PTCAttesterRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, fmt.Errorf("no post consensus roots for ptc attestation")
}

// executeDuty steps:
//  1. observe the slot's payload attestation data from the local beacon node; a zero BeaconBlockRoot
//     means no block was seen — abstain, freezing, signing and broadcasting nothing (SIP #94 §3)
//  2. freeze the observation, sign it under DomainPTCAttester and broadcast the partial signature
//  3. once a quorum of operators converged on the same data, reconstruct and submit the message
func (r *PTCAttesterRunner) executeDuty(duty types.Duty) error {
	slot := duty.DutySlot()

	data, err := r.beacon.GetPayloadAttestationData(slot)
	if err != nil {
		return errors.Wrap(err, "failed to get payload attestation data")
	}
	if data.BeaconBlockRoot == (phase0.Root{}) {
		// Abstain: the duty stays running with nothing frozen, so incoming peer partials are
		// rejected; a re-triggered duty re-observes from scratch (StartNewDuty clears the freeze).
		return nil
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
