package dkg

import (
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operator *Operator
	// InitMsg holds the init method which started this runner
	InitMsg *Init
	// Identifier unique for DKG session
	Identifier RequestID
	// DepositDataRoot is the signing root for the deposit data
	DepositDataRoot []byte
	// DepositDataSignatures holds partial sigs on deposit data
	DepositDataSignatures map[types.OperatorID]*PartialDepositData
	// PartialSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID][]byte
	// OutputMsgs holds all output messages received
	OutputMsgs map[types.OperatorID]*SignedOutput

	I uint16

	keygenSubProtocol Protocol
	signSubProtocol   Protocol
	config            *Config
}

func (r *Runner) Start() error {
	data, err := r.InitMsg.Encode()
	if err != nil {
		return err
	}
	outgoing, err := r.keygenSubProtocol.ProcessMsg(&base.Message{
		Header: &base.MessageHeader{
			SessionId: r.Identifier[:],
			MsgType:   int32(InitMsgType),
			Sender:    0,
			Receiver:  0,
		},
		Data: data,
	})
	if err != nil {
		return err
	}
	for _, message := range outgoing {
		err = r.signAndBroadcast(&message)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessMsg processes a DKG signed message and returns true and signed output if finished
func (r *Runner) ProcessMsg(msg *SignedMessage) (bool, map[types.OperatorID]*SignedOutput, error) {
	// TODO - validate message

	switch msg.Message.Header.MsgType {
	case int32(ProtocolMsgType):
		outgoing, err := r.keygenSubProtocol.ProcessMsg(msg.Message)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to process dkg msg")
		}
		err = r.broadcastMessages(outgoing, ProtocolMsgType)
		if err != nil {
			return false, nil, err
		}

		if hasOutput(outgoing, KeygenOutputType) {
			outputMsg := outgoing[len(outgoing)-1]
			keygenOutput := &KeygenOutput{}
			err = keygenOutput.Decode(outputMsg.Data)
			if err != nil {
				return false, nil, err
			}
			r.signSubProtocol = NewSignDepositData(r.InitMsg, keygenOutput, ProtocolConfig{
				Identifier:    r.Identifier,
				Operator:      r.Operator,
				BeaconNetwork: r.config.BeaconNetwork,
				Signer:        r.config.Signer,
			})
			outgoing1, err := r.signSubProtocol.Start()
			if err != nil {
				return false, nil, err
			}
			err = r.broadcastMessages(outgoing1, ProtocolMsgType)
		}
	case int32(DepositDataMsgType):
		outgoing, err := r.signSubProtocol.ProcessMsg(msg.Message)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to partial sig msg")
		}
		if hasOutput(outgoing, PartialOutputMsgType) {
			return true, nil, err
		}

		/*
				// TODO: Do we need to aggregate the signed outputs.
			case DepositDataMsgType:
				depSig := &PartialDepositData{}
				if err := depSig.Decode(msg.Message.Data); err != nil {
					return false, nil, errors.Wrap(err, "could not decode PartialDepositData")
				}

				if err := r.validateDepositDataSig(depSig); err != nil {
					return false, nil, errors.Wrap(err, "PartialDepositData invalid")
				}

				r.DepositDataSignatures[msg.Signer] = depSig
				if len(r.DepositDataSignatures) == int(r.InitMsg.Threshold) {
					// reconstruct deposit data sig
					depositSig, err := r.reconstructDepositDataSignature()
					if err != nil {
						return false, nil, errors.Wrap(err, "could not reconstruct deposit data sig")
					}

					// encrypt Operator's share
					encryptedShare, err := r.config.Signer.Encrypt(r.Operator.EncryptionPubKey, r.KeyGenOutput.Share.Serialize())
					if err != nil {
						return false, nil, errors.Wrap(err, "could not encrypt share")
					}

					ret, err := r.generateSignedOutput(&Output{
						RequestID:            r.Identifier,
						EncryptedShare:       encryptedShare,
						SharePubKey:          r.KeyGenOutput.Share.GetPublicKey().Serialize(),
						ValidatorPubKey:      r.KeyGenOutput.ValidatorPK,
						DepositDataSignature: depositSig,
					})
					if err != nil {
						return false, nil, errors.Wrap(err, "could not generate dkg SignedOutput")
					}

					if err := r.signAndBroadcastMsg(ret, OutputMsgType); err != nil {
						return false, nil, errors.Wrap(err, "could not broadcast SignedOutput")
					}
					return false, nil, nil
				} */
	case int32(OutputMsgType):
		output := &SignedOutput{}
		if err := output.Decode(msg.Message.Data); err != nil {
			return false, nil, errors.Wrap(err, "could not decode SignedOutput")
		}

		if err := r.validateSignedOutput(output); err != nil {
			return false, nil, errors.Wrap(err, "signed output invali")
		}

		r.OutputMsgs[msg.Signer] = output
		if len(r.OutputMsgs) == int(r.InitMsg.Threshold) {
			return true, r.OutputMsgs, nil
		}
		return false, nil, nil
	default:
		return false, nil, errors.New("msg type invalid")
	}

	return false, nil, nil
}

func (r *Runner) generateSignedOutput(o *Output) (*SignedOutput, error) {
	sig, err := r.config.Signer.SignDKGOutput(o, r.Operator.ETHAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign output")
	}

	return &SignedOutput{
		Data:      o,
		Signer:    r.Operator.OperatorID,
		Signature: sig,
	}, nil
}

func (r *Runner) broadcastMessages(msgs []base.Message, msgType MsgType) error {
	for _, message := range msgs {
		if message.Header.MsgType == int32(msgType) {
			err := r.signAndBroadcast(&message)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Runner) signAndBroadcast(msg *base.Message) error {
	sig, err := r.config.Signer.SignDKGOutput(msg, r.Operator.ETHAddress)
	if err != nil {
		return err
	}
	r.config.Network.Broadcast(&SignedMessage{
		Message:   msg,
		Signer:    r.Operator.OperatorID,
		Signature: sig,
	})
	return nil
}

func hasOutput(msgs []base.Message, msgType MsgType) bool {
	return msgs != nil && len(msgs) > 0 && msgs[len(msgs)-1].Header.MsgType == int32(msgType)
}

func (r *Runner) validateSignedOutput(msg *SignedOutput) error {
	panic("implement")
}
