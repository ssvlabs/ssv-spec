package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

type ValidatorCommitBoost struct {
	PreconfRunners  PreconfRunners
	BeaconNetwork   types.BeaconNetwork
	Network         Network
	Beacon          BeaconNode
	PreconfSidecar  PreconfSidecar
	CommitteeMember *types.CommitteeMember
	Share           *types.Share
	Signer          types.BeaconSigner
	OperatorSigner  *types.OperatorSigner
}

func NewValidatorCommitBoost(
	beaconNetwork types.BeaconNetwork,
	network Network,
	beacon BeaconNode,
	preconfSidecar PreconfSidecar,
	committeeMember *types.CommitteeMember,
	share *types.Share,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	runners map[phase0.Root]Runner,
) *ValidatorCommitBoost {
	return &ValidatorCommitBoost{
		BeaconNetwork:   beaconNetwork,
		PreconfRunners:  runners,
		Network:         network,
		Beacon:          beacon,
		PreconfSidecar:  preconfSidecar,
		Share:           share,
		CommitteeMember: committeeMember,
		Signer:          signer,
		OperatorSigner:  operatorSigner,
	}
}

// StartDuty starts a preconf duty for a validator given a signing request
func (v *ValidatorCommitBoost) StartDuty(request types.PreconfRequest) error {
	_, exist := v.PreconfRunners[request.Root]
	if exist {
		return errors.Errorf("duty runner for request %s already exists", request.Root.String())
	}
	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	shareMap[v.Share.ValidatorIndex] = v.Share
	dutyRunner, err := NewPreconfRunner(v.BeaconNetwork, shareMap, v.Beacon, v.PreconfSidecar, v.Network, v.Signer, v.OperatorSigner)
	if err != nil {
		return errors.Wrap(err, "failed to create new preconf runner")
	}
	v.PreconfRunners[request.Root] = dutyRunner
	var signingDuty = types.ValidatorDuty{
		Type:           types.BNRolePreconfirmation,
		PubKey:         phase0.BLSPubKey(v.Share.ValidatorPubKey),
		Slot:           v.BeaconNetwork.EstimatedCurrentSlot(), //TODO: check if this is correct
		ValidatorIndex: v.Share.ValidatorIndex,
	}
	return dutyRunner.StartNewDuty(&signingDuty, v.CommitteeMember.GetQuorum())
}

// ProcessMessage processes Network Message of all types
func (v *ValidatorCommitBoost) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	// Validate message
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := types.Verify(signedSSVMessage, v.CommitteeMember.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	msg := signedSSVMessage.SSVMessage

	cbPartialSigMsg := &types.CBPartialSignature{}
	if err := cbPartialSigMsg.Decode(msg.GetData()); err != nil {
		return errors.Wrap(err, "could not get commit boost partial sig message from network message")
	}

	requestRoot := cbPartialSigMsg.RequestRoot

	// Get runner
	dutyRunner := v.PreconfRunners[requestRoot]
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID")
	}

	// Validate message for runner
	if err := v.validateMessage(msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	switch msg.GetType() {
	case types.SSVPartialSignatureMsgType:
		// Decode
		psigMsgs := &types.PartialSignatureMessages{}
		if err := psigMsgs.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		// Validate
		if len(signedSSVMessage.OperatorIDs) != 1 {
			return errors.New("PartialSignatureMessage has more than 1 signer")
		}

		if err := psigMsgs.ValidateForSigner(signedSSVMessage.OperatorIDs[0]); err != nil {
			return errors.Wrap(err, "invalid PartialSignatureMessages")
		}

		// Process
		if psigMsgs.Type == types.PostConsensusPartialSig {
			return dutyRunner.ProcessPostConsensus(psigMsgs)
		}
		return dutyRunner.ProcessPreConsensus(psigMsgs)
	default:
		return errors.New("unknown msg")
	}
}

func (v *ValidatorCommitBoost) validateMessage(msg *types.SSVMessage) error {
	if !v.Share.ValidatorPubKey.MessageIDBelongs(msg.GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}
