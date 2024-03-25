package validation

import (
	"context"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

// MsgValidatorFunc represents a message validator
type MsgValidatorFunc = func(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult

func MsgValidation(runner ssv.Runner) MsgValidatorFunc {
	return func(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
		ssvMsg, err := DecodePubsubMsg(msg)
		if err != nil {
			return pubsub.ValidationReject
		}
		if validateSSVMessage(runner, ssvMsg) != nil {
			return pubsub.ValidationReject
		}

		switch ssvMsg.GetType() {
		case types.SSVConsensusMsgType:
			if validateConsensusMsg(runner, ssvMsg.Data) != nil {
				return pubsub.ValidationReject
			}
		case types.SSVPartialSignatureMsgType:
			if validatePartialSigMsg(runner, ssvMsg.Data) != nil {
				return pubsub.ValidationReject
			}
		default:
			return pubsub.ValidationReject
		}

		return pubsub.ValidationAccept
	}
}

func DecodePubsubMsg(msg *pubsub.Message) (*types.SSVMessage, error) {
	byts := msg.GetData()
	ret := &types.SSVMessage{}
	if err := ret.Decode(byts); err != nil {
		return nil, err
	}
	return ret, nil
}

func validateSSVMessage(runner ssv.Runner, msg *types.SSVMessage) error {
	if !runner.GetBaseRunner().Share.ValidatorPubKey.MessageIDBelongs(msg.GetID()) {
		return errors.New("msg ID doesn't match validator ID")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}

func validateConsensusMsg(runner ssv.Runner, data []byte) error {
	signedMsg := &qbft.SignedMessage{}
	if err := signedMsg.Decode(data); err != nil {
		return err
	}

	contr := runner.GetBaseRunner().QBFTController

	if err := contr.BaseMsgValidation(signedMsg); err != nil {
		return err
	}

	/**
	Main controller processing flow
	_______________________________
	All decided msgs are processed the same, out of instance
	All valid future msgs are saved in a container and can trigger highest decided futuremsg
	All other msgs (not future or decided) are processed normally by an existing instance (if found)
	*/
	if qbft.IsDecidedMsg(contr.Share, signedMsg) {
		return qbft.ValidateDecided(contr.GetConfig(), signedMsg, contr.Share)
	} else if signedMsg.Message.Height > contr.Height {
		return validateFutureMsg(contr.GetConfig(), signedMsg, contr.Share.Committee)
	} else {
		if inst := contr.StoredInstances.FindInstance(signedMsg.Message.Height); inst != nil {
			return inst.BaseMsgValidation(signedMsg)
		}
		return errors.New("unknown instance")
	}
}

func validatePartialSigMsg(runner ssv.Runner, data []byte) error {
	signedMsg := &types.SignedPartialSignatureMessage{}
	if err := signedMsg.Decode(data); err != nil {
		return err
	}

	if signedMsg.Message.Type == types.PostConsensusPartialSig {
		return runner.GetBaseRunner().ValidatePostConsensusMsg(runner, signedMsg)
	}
	return runner.GetBaseRunner().ValidatePreConsensusMsg(runner, signedMsg)
}

func validateFutureMsg(
	config qbft.IConfig,
	msg *qbft.SignedMessage,
	operators []*types.Operator,
) error {
	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if len(msg.GetSigners()) != 1 {
		return errors.New("allows 1 signer")
	}

	// verify signature
	if err := msg.Signature.VerifyByOperators(msg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}
