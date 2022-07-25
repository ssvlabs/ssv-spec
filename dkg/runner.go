package dkg

import (
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
	// PartialSignatures holds partial sigs on deposit data
	PartialSignatures map[types.OperatorID]*dkgtypes.ParsedPartialSigMessage
	// OutputMsgs holds all output messages received
	OutputMsgs map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage

	KeygenSubProtocol dkgtypes.Protocol
	keygenOutput      *dkgtypes.LocalKeyShare
	signOutput        *dkgtypes.SignedDepositDataMsgBody
	Config            *dkgtypes.Config
}

func (r *Runner) Start() error {
	outgoing, err := r.KeygenSubProtocol.Start()
	if err != nil {
		return err
	}
	err = r.broadcastMessages(outgoing, dkgtypes.ProtocolMsgType)
	if err != nil {
		return err
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
	case int32(dkgtypes.PartialSingatureMsgType):

		if msg.Header.MsgType != int32(dkgtypes.PartialSingatureMsgType) {
			return false, nil, errors.New("invalid message type")
		}

		err := r.handlePartialSigMessage(msg)
		if err != nil {
			return false, nil, err
		}

		if err != nil {
			return false, nil, errors.Wrap(err, "failed to partial sig msg")
		}

		if r.signOutput != nil {
			sig, err := r.Config.Signer.SignDKGOutput(r.signOutput, r.Operator.ETHAddress)
			if err != nil {
				return false, nil, err
			}
			r.signOutput.OperatorSignature = sig
			out := &dkgtypes.ParsedSignedDepositDataMessage{
				Header: &dkgtypes.MessageHeader{
					SessionId: r.Identifier[:],
					MsgType:   int32(dkgtypes.SignedDepositDataMsgType),
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
			r.broadcastMessages([]dkgtypes.Message{*base}, dkgtypes.SignedDepositDataMsgType)
			return false, nil, nil
		}
	case int32(dkgtypes.SignedDepositDataMsgType):
		output := &dkgtypes.ParsedSignedDepositDataMessage{}
		if err := output.FromBase(msg); err != nil {
			return false, nil, errors.Wrap(err, "could not decode SignedOutput")
		}
		if output.Header.RequestID() != r.Identifier {
			return false, nil, errors.New("request id mismatch")
		}
		r.OutputMsgs[types.OperatorID(msg.Header.Sender)] = output
		if len(r.OutputMsgs) == len(r.InitMsg.OperatorIds) {
			for _, message := range r.OutputMsgs {
				if message.Header.RequestID() != r.Identifier {
					return true, r.OutputMsgs, errors.New("one of more messages have mismatched request id")
				}
			}
			return true, r.OutputMsgs, nil
		}
		return false, nil, nil
	default:
		return false, nil, errors.New("msg type invalid")
	}

	return false, nil, nil
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
