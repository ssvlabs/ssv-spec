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
	Operator          *types.Operator
	Share             *types.Share
	Signer            types.BeaconSigner
	OperatorSigner    types.OperatorSigner
	SignatureVerifier types.SignatureVerifier
}

func NewValidator(
	network Network,
	beacon BeaconNode,
	operator *types.Operator,
	share *types.Share,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
	runners map[RunnerRole]Runner,
	signatureVerifier types.SignatureVerifier,
) *Validator {
	return &Validator{
		DutyRunners:       runners,
		Network:           network,
		Beacon:            beacon,
		Share:             share,
		Operator:          operator,
		Signer:            signer,
		OperatorSigner:    operatorSigner,
		SignatureVerifier: signatureVerifier,
	}
}

// StartDuty starts a duty for the validator
func (v *Validator) StartDuty(duty *types.BeaconDuty) error {
	role := types.MapDutyToRunnerRole(duty.Type)
	dutyRunner := v.DutyRunners[role]
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

	// Verify SignedSSVMessage's signature
	if err := v.SignatureVerifier.Verify(signedSSVMessage, v.Operator.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	msg := signedSSVMessage.SSVMessage

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
		qbftMsg := &qbft.Message{}
		if err := qbftMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}

		if err := qbftMsg.Validate(); err != nil {
			return errors.Wrap(err, "invalid qbft Message")
		}

		// Process
		return dutyRunner.ProcessConsensus(signedSSVMessage)
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

func (v *Validator) validateMessage(runner Runner, msg *types.SSVMessage) error {
	if !v.Share.ValidatorPubKey.MessageIDBelongs(msg.GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}
