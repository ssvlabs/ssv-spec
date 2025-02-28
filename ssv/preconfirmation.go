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
}

func NewPreconfRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	preconf PreconfSidecar,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) (Runner, error) {

	if len(share) != 1 {
		return nil, errors.New("must have one share")
	}

	return &PreconfRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType:     types.RolePreconfirmation,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:         beacon,
		preconf:        preconf,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		valCheck:       valCheck,
	}, nil
}

func (r *PreconfRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	return r.BaseRunner.baseStartNewDuty(r, duty, quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *PreconfRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *PreconfRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	return errors.New("no pre consensus phase for preconf runner")
}

func (r *PreconfRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg, &types.ValidatorConsensusData{})
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	// specific duty sig
	cd := decidedValue.(*types.ValidatorConsensusData)
	requestRoot, err := cd.GetPreconfRequest()
	if err != nil {
		return errors.Wrap(err, "could not get aggregate and proof")
	}

	msg, err := r.BaseRunner.signBeaconObject(r, r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty), requestRoot,
		cd.Duty.Slot,
		types.DomainCommitBoost)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     cd.Duty.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

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

func (r *PreconfRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
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

		cd := &types.ValidatorConsensusData{}
		err = cd.Decode(r.GetState().DecidedValue)
		if err != nil {
			return errors.Wrap(err, "could not create consensus data")
		}
		preconfRequest, err := cd.GetPreconfRequest()
		if err != nil {
			return errors.Wrap(err, "could not get preconf request root")
		}

		if err := r.GetPreconfSidecar().SubmitCommitment(preconfRequest.Root, specSig); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed aggregate")
		}
	}
	r.GetState().Finished = true
	return nil
}

func (r *PreconfRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, errors.New("no post consensus roots for preconfirmation")
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *PreconfRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	cd := &types.ValidatorConsensusData{}
	err := cd.Decode(r.GetState().DecidedValue)
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not create consensus data")
	}

	preconfRequest, err := cd.GetPreconfRequest()
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not get preconf request root")
	}
	return []ssz.HashRoot{preconfRequest}, phase0.DomainType{}, nil
}

func (r *PreconfRunner) executeDuty(duty types.Duty) error {
	request, err := r.GetPreconfSidecar().GetNewRequest()
	if err != nil {
		return errors.Wrap(err, "failed to get attestation data")
	}
	if err := r.BaseRunner.decide(r, duty.DutySlot(), &request); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
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
	return nil
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
