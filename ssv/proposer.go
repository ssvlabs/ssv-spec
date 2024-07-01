package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type ProposerRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewProposerRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) Runner {
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
	}
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
	obj, ver, err := r.GetBeaconNode().GetBeaconBlock(duty.Slot, r.GetShare().Graffiti, fullSig)
	if err != nil {
		return errors.Wrap(err, "failed to get Beacon block")
	}

	byts, err := obj.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal beacon block")
	}

	input := &types.ValidatorConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, input.Duty.DutySlot(), input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}

	return nil
}

func (r *ProposerRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	// specific duty sig
	var blkToSign ssz.HashRoot

	cd := decidedValue.(*types.ValidatorConsensusData)
	if r.decidedBlindedBlock() {
		_, blkToSign, err = cd.GetBlindedBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get blinded block data")
		}
	} else {
		_, blkToSign, err = cd.GetBlockData()
		if err != nil {
			return errors.Wrap(err, "could not get block data")
		}
	}

	msg, err := r.BaseRunner.signBeaconObject(r, r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty), blkToSign,
		cd.Duty.Slot,
		types.DomainProposer)
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

func (r *ProposerRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
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

		validatorConsensusData := &types.ValidatorConsensusData{}
		err = validatorConsensusData.Decode(r.GetState().DecidedValue)
		if err != nil {
			return errors.Wrap(err, "could not create consensus data")
		}
		if r.decidedBlindedBlock() {
			vBlindedBlk, _, err := validatorConsensusData.GetBlindedBlockData()
			if err != nil {
				return errors.Wrap(err, "could not get blinded block")
			}

			if err := r.GetBeaconNode().SubmitBlindedBeaconBlock(vBlindedBlk, specSig); err != nil {
				return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed blinded Beacon block")
			}
		} else {
			vBlk, _, err := validatorConsensusData.GetBlockData()
			if err != nil {
				return errors.Wrap(err, "could not get block")
			}

			if err := r.GetBeaconNode().SubmitBeaconBlock(vBlk, specSig); err != nil {
				return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed Beacon block")
			}
		}
	}
	r.GetState().Finished = true
	return nil
}

// decidedBlindedBlock returns true if decided value has a blinded block, false if regular block
// WARNING!! should be called after decided only
func (r *ProposerRunner) decidedBlindedBlock() bool {
	validatorConsensusData := &types.ValidatorConsensusData{}
	err := validatorConsensusData.Decode(r.GetState().DecidedValue)
	if err != nil {
		return false
	}
	_, _, err = validatorConsensusData.GetBlindedBlockData()
	return err == nil
}

func (r *ProposerRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.GetState().StartingDuty.DutySlot())
	return []ssz.HashRoot{types.SSZUint64(epoch)}, types.DomainRandao, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ProposerRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	validatorConsensusData := &types.ValidatorConsensusData{}
	err := validatorConsensusData.Decode(r.GetState().DecidedValue)
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not create consensus data")
	}
	if r.decidedBlindedBlock() {
		_, data, err := validatorConsensusData.GetBlindedBlockData()
		if err != nil {
			return nil, phase0.DomainType{}, errors.Wrap(err, "could not get blinded block data")
		}
		return []ssz.HashRoot{data}, types.DomainProposer, nil
	}

	_, data, err := validatorConsensusData.GetBlockData()
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not get block data")
	}
	return []ssz.HashRoot{data}, types.DomainProposer, nil
}

// executeDuty steps:
// 1) sign a partial randao sig and wait for 2f+1 partial sigs from peers
// 2) reconstruct randao and send GetBeaconBlock to BN
// 3) start consensus on duty + block data
// 4) Once consensus decides, sign partial block and broadcast
// 5) collect 2f+1 partial sigs, reconstruct and broadcast valid block sig to the BN
func (r *ProposerRunner) executeDuty(duty types.Duty) error {
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

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

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

func (r *ProposerRunner) GetOperatorSigner() types.OperatorSigner {
	return r.operatorSigner
}
