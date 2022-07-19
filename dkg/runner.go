package dkg

import (
	"encoding/json"
	"fmt"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operator *dkgtypes.Operator
	// InitMsg holds the init method which started this runner
	InitMsg *dkgtypes.Init
	// Identifier unique for DKG session
	Identifier dkgtypes.RequestID
	// DepositDataRoot is the signing root for the deposit data
	DepositDataRoot []byte
	// DepositDataSignatures holds partial sigs on deposit data
	DepositDataSignatures map[types.OperatorID]*dkgtypes.PartialDepositData
	// PartialSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID][]byte
	// OutputMsgs holds all output messages received
	OutputMsgs map[types.OperatorID]*dkgtypes.SignedOutput

	I uint64

	KeygenSubProtocol dkgtypes.Protocol
	signSubProtocol   dkgtypes.Protocol
	config            *dkgtypes.Config
}

func (r *Runner) Start() error {
	data, err := r.InitMsg.Encode()
	if err != nil {
		return err
	}
	outgoing, err := r.KeygenSubProtocol.ProcessMsg(&dkgtypes.Message{
		Header: &dkgtypes.MessageHeader{
			SessionId: r.Identifier[:],
			MsgType:   int32(dkgtypes.InitMsgType),
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
func (r *Runner) ProcessMsg(msg *dkgtypes.Message) (bool, map[types.OperatorID]*dkgtypes.SignedOutput, error) {
	// TODO - validate message

	switch msg.Header.MsgType {
	case int32(dkgtypes.ProtocolMsgType):
		outgoing, err := r.KeygenSubProtocol.ProcessMsg(msg)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to process dkg msg")
		}
		err = r.broadcastMessages(outgoing, dkgtypes.ProtocolMsgType)
		if err != nil {
			return false, nil, err
		}

		output, err := r.KeygenSubProtocol.Output()
		if err != nil {
			return false, nil, err
		}
		lks := &dkgtypes.LocalKeyShare{}
		json.Unmarshal(output, lks)
		jstr, err := json.Marshal(lks)
		fmt.Printf("output is %v\n", string(jstr))
		if hasOutput(outgoing, dkgtypes.KeygenOutputType) {
			outputMsg := outgoing[len(outgoing)-1]
			keygenOutput := &dkgtypes.KeygenOutput{}
			err = keygenOutput.Decode(outputMsg.Data)
			if err != nil {
				return false, nil, err
			}
			r.signSubProtocol = dkgtypes.NewSignDepositData(r.InitMsg, keygenOutput, dkgtypes.ProtocolConfig{
				Identifier:    r.Identifier,
				Operator:      r.Operator,
				BeaconNetwork: r.config.BeaconNetwork,
				Signer:        r.config.Signer,
			})
			outgoing1, err := r.signSubProtocol.Start()
			if err != nil {
				return false, nil, err
			}
			err = r.broadcastMessages(outgoing1, dkgtypes.ProtocolMsgType)
		}
	case int32(dkgtypes.DepositDataMsgType):
		outgoing, err := r.signSubProtocol.ProcessMsg(msg)
		fmt.Printf("outgoing is %v\n", outgoing)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to partial sig msg")
		}
		if hasOutput(outgoing, dkgtypes.PartialOutputMsgType) {
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
	case int32(dkgtypes.OutputMsgType):
		output := &dkgtypes.SignedOutput{}
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

func (r *Runner) generateSignedOutput(o *dkgtypes.Output) (*dkgtypes.SignedOutput, error) {
	sig, err := r.config.Signer.SignDKGOutput(o, r.Operator.ETHAddress)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign output")
	}

	return &dkgtypes.SignedOutput{
		Data:      o,
		Signer:    r.Operator.OperatorID,
		Signature: sig,
	}, nil
}

func (r *Runner) broadcastMessages(msgs []dkgtypes.Message, msgType dkgtypes.MsgType) error {
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

func (r *Runner) signAndBroadcast(msg *dkgtypes.Message) error {
	sig, err := r.config.Signer.SignDKGOutput(msg, r.Operator.ETHAddress)
	if err != nil {
		return err
	}
	r.config.Network.Broadcast(&dkgtypes.SignedMessage{
		Message:   msg,
		Signer:    r.Operator.OperatorID,
		Signature: sig,
	})
	return nil
}

func hasOutput(msgs []dkgtypes.Message, msgType dkgtypes.MsgType) bool {
	return msgs != nil && len(msgs) > 0 && msgs[len(msgs)-1].Header.MsgType == int32(msgType)
}

func (r *Runner) validateSignedOutput(msg *dkgtypes.SignedOutput) error {
	panic("implement")
}
