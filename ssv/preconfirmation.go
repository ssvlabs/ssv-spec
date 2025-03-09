package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type PreconfRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	preconf        PreconfSidecar
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF

	requestRoot phase0.Root
}

func NewPreconfRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	preconf PreconfSidecar,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
) (Runner, error) {

	if len(share) != 1 {
		return nil, errors.New("must have one share")
	}

	return &PreconfRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RolePreconfirmation,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},

		beacon:         beacon,
		preconf:        preconf,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		valCheck:       valCheck,
		requestRoot:    phase0.Root{},
	}, nil
}

func (r *PreconfRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	if err := r.ShouldProcessDuty(duty); err != nil {
		return errors.Wrap(err, "can't start duty")
	}

	r.BaseRunner.baseSetupForNewDuty(duty, quorum)
	return r.executeDuty(duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *PreconfRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *PreconfRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing preconfirmation message")
	}

	if !quorum {
		return nil
	}

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

	if err := r.GetPreconfSidecar().SubmitCommitment(r.requestRoot, specSig); err != nil {
		return errors.Wrap(err, "could not submit to commitment to sidecar")
	}

	r.GetState().Finished = true
	r.requestRoot = phase0.Root{}
	return nil
}

func (r *PreconfRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return errors.New("no consensus phase for preconfirmation")
}

func (r *PreconfRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return errors.New("no post consensus phase for preconfirmation")
}

func (r *PreconfRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	if r.BaseRunner.State == nil || r.BaseRunner.State.StartingDuty == nil || r.requestRoot == (phase0.Root{}) {
		return nil, types.DomainError, errors.New("no running duty or preconf request")
	}
	preconfRequest := types.PreconfRequest{
		Root: r.requestRoot,
	}
	return []ssz.HashRoot{&preconfRequest}, types.DomainCommitBoost, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *PreconfRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, errors.New("no post consensus roots for preconfirmation")
}

func (r *PreconfRunner) executeDuty(duty types.Duty) error {
	request, err := r.GetPreconfSidecar().GetNewRequest()
	if err != nil {
		return errors.Wrap(err, "failed to get preconf request")
	}

	r.requestRoot = request.Root

	msg, err := r.BaseRunner.signBeaconObject(r, r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty), &request,
		duty.DutySlot(),
		types.DomainCommitBoost)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}
	preConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PreconfPartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := preConsensusMsg.Encode()
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

// override ShouldProcessDuty to allow multiple duties in the same slot
func (r *PreconfRunner) ShouldProcessDuty(duty types.Duty) error {
	if r.GetState() != nil && r.GetState().StartingDuty.DutySlot() > duty.DutySlot() {
		return errors.Errorf("duty for slot %d already passed. Current height is %d", duty.DutySlot(),
			r.BaseRunner.QBFTController.Height)
	}
	// multiple preconf duties are allowed in the same slot, but only one can be running at a time
	if r.requestRoot != (phase0.Root{}) || r.HasRunningDuty() {
		return errors.Errorf("has a running duty, try after the current duty finishes")
	}
	return nil
}

func (r *PreconfRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *PreconfRunner) GetNetwork() Network {
	return r.network
}

func (r *PreconfRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *PreconfRunner) GetPreconfSidecar() PreconfSidecar {
	return r.preconf
}

func (r *PreconfRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *PreconfRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *PreconfRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *PreconfRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *PreconfRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}
