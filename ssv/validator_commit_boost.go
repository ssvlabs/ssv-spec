package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

type ValidatorCommitBoost struct {
	CBSigningRunners CBSigningRunners
	BeaconNetwork    types.BeaconNetwork
	Network          Network
	Beacon           BeaconNode
	CommitteeMember  *types.CommitteeMember
	Share            *types.Share
	Signer           types.BeaconSigner
	OperatorSigner   *types.OperatorSigner
}

func NewValidatorCommitBoost(
	beaconNetwork types.BeaconNetwork,
	network Network,
	beacon BeaconNode,
	committeeMember *types.CommitteeMember,
	share *types.Share,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
) *ValidatorCommitBoost {
	return &ValidatorCommitBoost{
		BeaconNetwork:    beaconNetwork,
		CBSigningRunners: make(CBSigningRunners),
		Network:          network,
		Beacon:           beacon,
		Share:            share,
		CommitteeMember:  committeeMember,
		Signer:           signer,
		OperatorSigner:   operatorSigner,
	}
}

func (v *ValidatorCommitBoost) HandleRequestSignature(keyType string, pubkey types.ValidatorPK, objectRoot phase0.Root) (phase0.BLSSignature, error) {
	// Proxy key is not supported currently
	if keyType != "consensus" {
		return phase0.BLSSignature{}, errors.New("invalid key type")
	}

	if pubkey != v.Share.ValidatorPubKey {
		return phase0.BLSSignature{}, errors.New("invalid pubkey")
	}

	var signingDuty = types.CBSigningDuty{
		Request: types.CBSigningRequest{
			Root: objectRoot,
		},
		Duty: types.ValidatorDuty{
			Slot:           v.BeaconNetwork.EstimatedCurrentSlot(),
			ValidatorIndex: v.Share.ValidatorIndex,
		},
	}

	err := v.StartDuty(signingDuty)
	if err != nil {
		return phase0.BLSSignature{}, errors.Wrap(err, "failed to start duty")
	}

	dutyRunner, exist := v.CBSigningRunners[objectRoot]
	if !exist {
		return phase0.BLSSignature{}, errors.Errorf("could not get duty runner for request %s", objectRoot.String())
	}
	sig := dutyRunner.GetSignature()

	return sig, nil
}

// StartDuty starts a cb signing duty for a validator given a signing request
func (v *ValidatorCommitBoost) StartDuty(duty types.CBSigningDuty) error {
	_, exist := v.CBSigningRunners[duty.Request.Root]
	if exist {
		return errors.Errorf("duty runner for request %s already exists", duty.Request.Root.String())
	}
	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	shareMap[v.Share.ValidatorIndex] = v.Share
	dutyRunner, err := NewCBSigningRunner(v.BeaconNetwork, shareMap, v.Beacon, v.Network, v.Signer, v.OperatorSigner)
	if err != nil {
		return errors.Wrap(err, "failed to create new commit-boost signing runner")
	}
	v.CBSigningRunners[duty.Request.Root] = dutyRunner
	return dutyRunner.StartNewDuty(duty, v.CommitteeMember.GetQuorum())
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

	cbPartialSigMsg := &types.CBPartialSignatures{}
	if err := cbPartialSigMsg.Decode(msg.GetData()); err != nil {
		return errors.Wrap(err, "could not get commit boost partial sig message from network message")
	}

	requestRoot := cbPartialSigMsg.RequestRoot

	// Get runner
	dutyRunner, exist := v.CBSigningRunners[requestRoot]
	if !exist {
		return errors.Errorf("could not get duty runner for request %s", requestRoot.String())
	}

	// Validate message for runner
	if err := v.validateMessage(msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	switch msg.GetType() {
	case types.CommitBoostPartialSignatureMsgType:
		psigMsgs := cbPartialSigMsg.PartialSig
		// Validate
		if len(signedSSVMessage.OperatorIDs) != 1 {
			return errors.New("PartialSignatureMessage has more than 1 signer")
		}

		if err := cbPartialSigMsg.PartialSig.ValidateForSigner(signedSSVMessage.OperatorIDs[0]); err != nil {
			return errors.Wrap(err, "invalid PartialSignatureMessages")
		}

		// Process
		if cbPartialSigMsg.PartialSig.Type == types.PostConsensusPartialSig {
			return dutyRunner.ProcessPostConsensus(&psigMsgs)
		}
		return dutyRunner.ProcessPreConsensus(&psigMsgs)
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
