package dkg

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/dkg/base"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operator *base.Operator
	// InitMsg holds the init method which started this runner
	InitMsg *base.Init
	// Identifier unique for DKG session
	Identifier base.RequestID
	// DepositDataRoot is the signing root for the deposit data
	DepositDataRoot []byte
	// DepositDataSignatures holds partial sigs on deposit data
	DepositDataSignatures map[types.OperatorID]*base.PartialDepositData
	// PartialSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID][]byte
	// OutputMsgs holds all output messages received
	OutputMsgs map[types.OperatorID]*base.SignedOutput

	I uint64

	KeygenSubProtocol base.Protocol
	signSubProtocol   base.Protocol
	config            *base.Config
}

func (r *Runner) Start() error {
	data, err := r.InitMsg.Encode()
	if err != nil {
		return err
	}
	outgoing, err := r.KeygenSubProtocol.ProcessMsg(&base.Message{
		Header: &base.MessageHeader{
			SessionId: r.Identifier[:],
			MsgType:   int32(base.InitMsgType),
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
func (r *Runner) ProcessMsg(msg *base.Message) (bool, map[types.OperatorID]*base.SignedOutput, error) {
	// TODO - validate message

	switch msg.Header.MsgType {
	case int32(base.ProtocolMsgType):
		outgoing, err := r.KeygenSubProtocol.ProcessMsg(msg)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to process dkg msg")
		}
		err = r.broadcastMessages(outgoing, base.ProtocolMsgType)
		if err != nil {
			return false, nil, err
		}

		output, err := r.KeygenSubProtocol.Output()
		if err != nil {
			return false, nil, err
		}
		lks := &base.LocalKeyShare{}
		json.Unmarshal(output, lks)
		jstr, err := json.Marshal(lks)
		fmt.Printf("output is %v\n", string(jstr))
		if hasOutput(outgoing, base.KeygenOutputType) {
			outputMsg := outgoing[len(outgoing)-1]
			keygenOutput := &base.KeygenOutput{}
			err = keygenOutput.Decode(outputMsg.Data)
			if err != nil {
				return false, nil, err
			}
			r.signSubProtocol = base.NewSignDepositData(r.InitMsg, keygenOutput, base.ProtocolConfig{
				Identifier:    r.Identifier,
				Operator:      r.Operator,
				BeaconNetwork: r.config.BeaconNetwork,
				Signer:        r.config.Signer,
			})
			outgoing1, err := r.signSubProtocol.Start()
			if err != nil {
				return false, nil, err
			}
			err = r.broadcastMessages(outgoing1, base.ProtocolMsgType)
		}
	case int32(base.DepositDataMsgType):
		outgoing, err := r.signSubProtocol.ProcessMsg(msg)
		fmt.Printf("outgoing is %v\n", outgoing)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to partial sig msg")
		}
		if hasOutput(outgoing, base.PartialOutputMsgType) {
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
	case int32(base.OutputMsgType):
		output := &base.SignedOutput{}
		if err := output.Decode(msg.Data); err != nil {
			return false, nil, errors.Wrap(err, "could not decode SignedOutput")
		}

		if err := r.validateSignedOutput(output); err != nil {
			return false, nil, errors.Wrap(err, "signed output invali")
		}

		r.OutputMsgs[types.OperatorID(msg.Header.Sender)] = output
		if len(r.OutputMsgs) == int(r.InitMsg.Threshold) {
			return true, r.OutputMsgs, nil
		}
		return false, nil, nil
	default:
		return false, nil, errors.New("msg type invalid")
	}

	return false, nil, nil
}

func (r *Runner) generateSignedOutput(o *base.Output) (*base.SignedOutput, error) {
	sig, err := r.config.Signer.SignDKGOutput(o, r.Operator.ETHAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign output")
	}

	return &base.SignedOutput{
		Data:      o,
		Signer:    r.Operator.OperatorID,
		Signature: sig,
	}, nil
}

func (r *Runner) broadcastMessages(msgs []base.Message, msgType base.MsgType) error {
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
	r.config.Network.Broadcast(&base.SignedMessage{
		Message:   msg,
		Signer:    r.Operator.OperatorID,
		Signature: sig,
	})
	return nil
}

func hasOutput(msgs []base.Message, msgType base.MsgType) bool {
	return msgs != nil && len(msgs) > 0 && msgs[len(msgs)-1].Header.MsgType == int32(msgType)
}

func (r *Runner) validateSignedOutput(msg *base.SignedOutput) error {
	panic("implement")
}
