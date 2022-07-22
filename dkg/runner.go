package dkg

import (
	"github.com/bloxapp/ssv-spec/dkg/sign"
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
	DepositDataSignatures map[types.OperatorID]*dkgtypes.PartialSigMsgBody
	// PartialSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID][]byte
	// OutputMsgs holds all output messages received
	OutputMsgs map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage

	KeygenSubProtocol dkgtypes.Protocol
	SignSubProtocol   dkgtypes.Protocol
	keygenOutput      *dkgtypes.LocalKeyShare
	signOutput        *dkgtypes.SignedDepositDataMsgBody
	Config            *dkgtypes.Config
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
func (r *Runner) ProcessMsg(msg *dkgtypes.Message) (bool, map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage, error) {
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

		if output != nil {
			r.keygenOutput = &dkgtypes.LocalKeyShare{}
			if err = r.keygenOutput.Decode(output); err != nil {
				return false, nil, err
			}
			if err = r.startSigning(); err != nil {
				return false, nil, err
			}
		}
	case int32(dkgtypes.DepositDataMsgType):
		_, err := r.SignSubProtocol.ProcessMsg(msg)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to partial sig msg")
		}
		output, err := r.SignSubProtocol.Output()
		if err != nil {
			return false, nil, err
		}
		if output != nil {
			// TODO do we need to store the output?
			r.signOutput = &dkgtypes.SignedDepositDataMsgBody{}
			err = r.signOutput.Decode(output)
			sig, err := r.Config.Signer.SignDKGOutput(r.signOutput, r.Operator.ETHAddress)
			if err != nil {
				return false, nil, err
			}
			r.signOutput.OperatorSignature = sig
			out := &dkgtypes.ParsedSignedDepositDataMessage{
				Header: &dkgtypes.MessageHeader{
					SessionId: r.Identifier[:],
					MsgType:   int32(dkgtypes.OutputMsgType),
					Sender:    uint64(r.Operator.OperatorID),
					Receiver:  0,
				},
				Body:      r.signOutput,
				Signature: nil,
			}
			base, err := out.ToBase()
			if err != nil {
				return false, nil, err
			}
			r.OutputMsgs[r.Operator.OperatorID] = out
			r.broadcastMessages([]dkgtypes.Message{*base}, dkgtypes.OutputMsgType)
			return false, nil, nil
		}
		//if hasOutput(outgoing, dkgtypes.PartialOutputMsgType) {
		//	return true, nil, err
		//}

		/*
				// TODO: Do we need to aggregate the signed outputs.
			case DepositDataMsgType:
				depSig := &PartialSigMsgBody{}
				if err := depSig.Decode(msg.Message.Data); err != nil {
					return false, nil, errors.Wrap(err, "could not decode PartialSigMsgBody")
				}

				if err := r.validateDepositDataSig(depSig); err != nil {
					return false, nil, errors.Wrap(err, "PartialSigMsgBody invalid")
				}

				r.DepositDataSignatures[msg.Signer] = depSig
				if len(r.DepositDataSignatures) == int(r.InitMsg.Threshold) {
					// reconstruct deposit data sig
					depositSig, err := r.reconstructDepositDataSignature()
					if err != nil {
						return false, nil, errors.Wrap(err, "could not reconstruct deposit data sig")
					}

					// encrypt Operator's share
					encryptedShare, err := r.Config.Signer.Encrypt(r.Operator.EncryptionPubKey, r.KeyGenOutput.Share.Serialize())
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
		output := &dkgtypes.ParsedSignedDepositDataMessage{}
		if err := output.FromBase(msg); err != nil {
			return false, nil, errors.Wrap(err, "could not decode SignedOutput")
		}

		if err := r.validateSignedOutput(output); err != nil {
			return false, nil, errors.Wrap(err, "signed output invalid")
		}
		r.OutputMsgs[types.OperatorID(msg.Header.Sender)] = output
		if len(r.OutputMsgs) == len(r.InitMsg.OperatorIDs) {
			return true, r.OutputMsgs, nil
		}
		return false, nil, nil
	default:
		return false, nil, errors.New("msg type invalid")
	}

	return false, nil, nil
}

func (r *Runner) startSigning() error {

	r.SignSubProtocol = sign.NewSignDepositData(r.InitMsg, r.keygenOutput, dkgtypes.ProtocolConfig{
		Identifier:    r.Identifier,
		Operator:      r.Operator,
		BeaconNetwork: r.Config.BeaconNetwork,
		Signer:        r.Config.Signer,
	})

	outgoing, err := r.SignSubProtocol.Start()
	if err != nil {
		return err
	}
	err = r.broadcastMessages(outgoing, dkgtypes.ProtocolMsgType)
	err = r.broadcastMessages(outgoing, dkgtypes.DepositDataMsgType)
	return nil
}

//func (r *Runner) generateSignedOutput(o *dkgtypes.Output) (*dkgtypes.SignedOutput, error) {
//	sig, err := r.Config.Signer.SignDKGOutput(o, r.Operator.ETHAddress)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not sign output")
//	}
//
//	return &dkgtypes.SignedOutput{
//		Data:      o,
//		Signer:    r.Operator.OperatorID,
//		Signature: sig,
//	}, nil
//}

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
	sig, err := r.Config.Signer.SignDKGOutput(msg, r.Operator.ETHAddress)
	if err != nil {
		return err
	}
	err = msg.SetSignature(sig)
	if err != nil {
		return err
	}
	r.Config.Network.Broadcast(msg)
	return nil
}

func hasOutput(msgs []dkgtypes.Message, msgType dkgtypes.MsgType) bool {
	return msgs != nil && len(msgs) > 0 && msgs[len(msgs)-1].Header.MsgType == int32(msgType)
}

func (r *Runner) validateSignedOutput(msg *dkgtypes.ParsedSignedDepositDataMessage) error {
	if msg == nil {
		return errors.New("msg is nil")
	}
	if !r.signOutput.SameDepositData(msg.Body) {
		return errors.New("deposit data doesn't match")
	}
	// TODO: Verify OperatorSignature
	return nil
}
