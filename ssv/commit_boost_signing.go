package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type CBSigningRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF

	requestRoot phase0.Root
	requestSig  chan phase0.BLSSignature
}

func NewCBSigningRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
) (*CBSigningRunner, error) {

	if len(share) != 1 {
		return nil, errors.New("must have one share")
	}

	return &CBSigningRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RoleCBSigning,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		requestRoot:    phase0.Root{},
		requestSig:     make(chan phase0.BLSSignature, 1),
	}, nil
}

func (r *CBSigningRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	r.BaseRunner.baseSetupForNewDuty(duty, quorum)
	return r.executeDuty(duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *CBSigningRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *CBSigningRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing commit-boost signing message")
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

	r.requestSig <- specSig

	r.GetState().Finished = true
	return nil
}

func (r *CBSigningRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return errors.New("no consensus phase for commit-boost signing")
}

func (r *CBSigningRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return errors.New("no post consensus phase for commit-boost signing")
}

func (r *CBSigningRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	if r.BaseRunner.State == nil || r.BaseRunner.State.StartingDuty == nil {
		return nil, types.DomainError, errors.New("no running duty")
	}
	if r.requestRoot == (phase0.Root{}) {
		return nil, types.DomainError, errors.New("no request root")
	}
	CBSigningRequest := types.CBSigningRequest{
		Root: r.requestRoot,
	}
	return []ssz.HashRoot{&CBSigningRequest}, types.DomainCommitBoost, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *CBSigningRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, errors.New("no post consensus roots for commit-boost signing")
}

func (r *CBSigningRunner) executeDuty(duty types.Duty) error {
	cbSigningDuty := types.CBSigningDuty{}
	if cb, ok := duty.(*types.CBSigningDuty); ok {
		cbSigningDuty = *cb
	} else if cb, ok := duty.(types.CBSigningDuty); ok {
		cbSigningDuty = cb
	} else {
		return errors.New("duty is not a CBSigningDuty")
	}
	request := cbSigningDuty.Request

	r.requestRoot = request.Root

	msg, err := r.BaseRunner.signBeaconObject(r, &cbSigningDuty.Duty, &request,
		duty.DutySlot(),
		types.DomainCommitBoost)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}

	preConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.CBSigningPartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{msg},
	}

	CBPreConsensusMsg := &types.CBPartialSignatures{
		RequestRoot: r.requestRoot,
		PartialSig:  *preConsensusMsg,
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := CBPreConsensusMsg.Encode()
	if err != nil {
		return err
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.CommitBoostPartialSignatureMsgType,
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

func (r *CBSigningRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *CBSigningRunner) GetNetwork() Network {
	return r.network
}

func (r *CBSigningRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *CBSigningRunner) GetShare() *types.Share {
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *CBSigningRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *CBSigningRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *CBSigningRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *CBSigningRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}

func (r *CBSigningRunner) GetSignature() phase0.BLSSignature {
	return <-r.requestSig
}
