package ssv

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Validator represents an SSV ETH consensus validator Share assigned, coordinates duty execution and more.
// Every validator has a validatorID which is validator's public key.
// Each validator has multiple DutyRunners, for each duty type.
type Validator struct {
	DutyRunners       DutyRunners
	Network           Network
	Beacon            BeaconNode
	Share             *types.Share
	Signer            types.KeyManager
	OperatorSigner    types.OperatorSigner
	SignatureVerifier types.SignatureVerifier
}

func NewValidator(
	network Network,
	beacon BeaconNode,
	share *types.Share,
	signer types.KeyManager,
	operatorSigner types.OperatorSigner,
	runners map[types.BeaconRole]Runner,
	signatureVerifier types.SignatureVerifier,
) *Validator {
	return &Validator{
		DutyRunners:       runners,
		Network:           network,
		Beacon:            beacon,
		Share:             share,
		Signer:            signer,
		OperatorSigner:    operatorSigner,
		SignatureVerifier: signatureVerifier,
	}
}

// StartDuty starts a duty for the validator
func (v *Validator) StartDuty(duty *types.Duty) error {
	dutyRunner := v.DutyRunners[duty.Type]
	if dutyRunner == nil {
		return errors.Errorf("duty type %s not supported", duty.Type.String())
	}
	return dutyRunner.StartNewDuty(duty)
}

// ProcessMessage processes Network Message of all types
func (v *Validator) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {

	// Validate message
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Decode the nested SSVMessage
	msg := &types.SSVMessage{}
	if err := msg.Decode(signedSSVMessage.Data); err != nil {
		return errors.Wrap(err, "could not decode data into an SSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := v.SignatureVerifier.Verify(signedSSVMessage, v.Share.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	// Get runner
	dutyRunner := v.DutyRunners.DutyRunnerForMsgID(msg.GetID())
	if dutyRunner == nil {
		return errors.Errorf("could not get duty runner for msg ID")
	}

	// Validate message for runner
	if err := v.validateMessage(dutyRunner, msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		// Decode
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}

		// Check signer consistency
		if !signedMsg.CommonSigners([]types.OperatorID{signedSSVMessage.OperatorID}) {
			return errors.New("SignedSSVMessage's signer not consistent with SignedMessage's signers")
		}

		// Process
		return dutyRunner.ProcessConsensus(signedMsg)
	case types.SSVPartialSignatureMsgType:
		// Decode
		signedMsg := &types.SignedPartialSignatureMessage{}
		if err := signedMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		// Check signer consistency
		if signedMsg.Signer != signedSSVMessage.OperatorID {
			return errors.New("SignedSSVMessage's signer not consistent with SignedPartialSignatureMessage's signer")
		}

		// Process
		if signedMsg.Message.Type == types.PostConsensusPartialSig {
			return dutyRunner.ProcessPostConsensus(signedMsg)
		}
		return dutyRunner.ProcessPreConsensus(signedMsg)
	default:
		return errors.New("unknown msg")
	}
}

func (v *Validator) validateMessage(runner Runner, msg *types.SSVMessage) error {
	if !v.Share.ValidatorPubKey.MessageIDBelongs(msg.GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}
