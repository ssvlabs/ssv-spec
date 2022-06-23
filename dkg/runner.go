package dkg

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operator *Operator
	// InitMsg holds the init method which started this runner
	InitMsg *Init
	// Identifier unique for DKG session
	Identifier RequestID
	// KeyGenOutput holds the protocol output once it finishes
	KeyGenOutput *KeyGenOutput
	// DepositDataRoot is the signing root for the deposit data
	DepositDataRoot []byte
	// DepositDataSignatures holds partial sigs on deposit data
	DepositDataSignatures map[types.OperatorID]*PartialDepositData
	// OutputMsgs holds all output messages received
	OutputMsgs map[types.OperatorID]*SignedOutput

	protocol KeyGenProtocol
	config   *Config
}

// ProcessMsg processes a DKG signed message and returns true and signed output if finished
func (r *Runner) ProcessMsg(msg *SignedMessage) (bool, map[types.OperatorID]*SignedOutput, error) {
	// TODO - validate message

	switch msg.Message.MsgType {
	case ProtocolMsgType:
		finished, o, err := r.protocol.ProcessMsg(msg)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to process dkg msg")
		}

		if finished {
			r.KeyGenOutput = o

			// generate deposit data
			root, _, err := types.GenerateETHDepositData(
				r.KeyGenOutput.ValidatorPK,
				r.InitMsg.WithdrawalCredentials,
				r.InitMsg.Fork,
				types.DomainDeposit,
			)
			if err != nil {
				return false, nil, errors.Wrap(err, "could not generate deposit data")
			}
			r.DepositDataRoot = root

			// sign
			sig, err := r.config.Signer.SignETHDepositRoot(root, r.Operator.ETHAddress)
			if err != nil {
				return false, nil, errors.Wrap(err, "could not sign deposit data")
			}

			// broadcast
			if err := r.signAndBroadcastMsg(&PartialDepositData{
				Signer:    r.Operator.OperatorID,
				Root:      r.DepositDataRoot,
				Signature: sig,
			}, DepositDataMsgType); err != nil {
				return false, nil, errors.Wrap(err, "could not broadcast partial deposit data")
			}
		}
		return false, nil, nil
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
		}
	case OutputMsgType:
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

func (r *Runner) signAndBroadcastMsg(msg types.Encoder, msgType MsgType) error {
	panic("implement")
}

func (r *Runner) reconstructDepositDataSignature() (types.Signature, error) {
	panic("implement")
}

func (r *Runner) validateSignedOutput(msg *SignedOutput) error {
	panic("implement")
}

func (r *Runner) validateDepositDataSig(msg *PartialDepositData) error {
	if !bytes.Equal(r.DepositDataRoot, msg.Root) {
		return errors.New("deposit data roots not equal")
	}

	// find operator and verify msg
	sharePK, found := r.KeyGenOutput.OperatorPubKeys[msg.Signer]
	if !found {
		return errors.New("signer not part of committee")
	}
	sig := &bls.Sign{}
	if err := sig.Deserialize(msg.Signature); err != nil {
		return errors.Wrap(err, "could not deserialize partial sig")
	}
	if !sig.VerifyByte(sharePK, r.DepositDataRoot) {
		return errors.New("partial deposit data sig invalid")
	}

	return nil
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
