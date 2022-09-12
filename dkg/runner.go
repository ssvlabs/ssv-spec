package dkg

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
		if r.DepositDataSignatures[r.Operator.OperatorID] != nil {
			return false, nil, errors.New("keygen has already completed")
		}
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
			sig := r.KeyGenOutput.Share.SignByte(root)

			// broadcast
			pdd := &PartialDepositData{
				Signer:    r.Operator.OperatorID,
				Root:      r.DepositDataRoot,
				Signature: sig.Serialize(),
			}
			if err := r.signAndBroadcastMsg(pdd, DepositDataMsgType); err != nil {
				return false, nil, errors.Wrap(err, "could not broadcast partial deposit data")
			}
			r.DepositDataSignatures[r.Operator.OperatorID] = pdd
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

		if found := r.DepositDataSignatures[msg.Signer]; found == nil {
			r.DepositDataSignatures[msg.Signer] = depSig
		} else if !bytes.Equal(found.Signature, msg.Signature) {
			return false, nil, errors.New("inconsistent partial signature received")
		}

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

			r.OutputMsgs[r.Operator.OperatorID] = ret
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
		// GLNOTE: Actually we need every operator to sign instead only the quorum!
		if len(r.OutputMsgs) == len(r.InitMsg.OperatorIDs) {
			return true, r.OutputMsgs, nil
		}
		return false, nil, nil
	default:
		return false, nil, errors.New("msg type invalid")
	}

	return false, nil, nil
}

func (r *Runner) signAndBroadcastMsg(msg types.Encoder, msgType MsgType) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}
	signedMessage := &SignedMessage{
		Message: &Message{
			MsgType:    msgType,
			Identifier: r.Identifier,
			Data:       data,
		},
		Signer:    r.Operator.OperatorID,
		Signature: nil,
	}
	// GLNOTE: Should we use SignDKGOutput?
	sig, err := r.config.Signer.SignDKGOutput(signedMessage, r.Operator.ETHAddress)
	if err != nil {
		return errors.Wrap(err, "failed to sign message")
	}
	signedMessage.Signature = sig
	if err = r.config.Network.BroadcastDKGMessage(signedMessage); err != nil {
		return errors.Wrap(err, "failed to broadcast message")
	}
	return nil
}

func (r *Runner) reconstructDepositDataSignature() (types.Signature, error) {
	sigBytes := map[types.OperatorID][]byte{}
	for id, d := range r.DepositDataSignatures {
		if err := r.validateDepositDataRoot(d); err != nil {
			return nil, errors.Wrap(err, "PartialDepositData invalid")
		}
		sigBytes[id] = d.Signature
	}

	sig, err := types.ReconstructSignatures(sigBytes)
	if err != nil {
		return nil, err
	}
	return sig.Serialize(), nil
}

func (r *Runner) validateSignedOutput(msg *SignedOutput) error {
	// TODO: Separate fields match and signature validation
	output := r.ownOutput()
	if output != nil {
		if output.Data.RequestID != msg.Data.RequestID {
			return errors.New("got mismatching RequestID")
		}
		if !bytes.Equal(output.Data.ValidatorPubKey, msg.Data.ValidatorPubKey) {
			return errors.New("got mismatching ValidatorPubKey")
		}
	}
	found, operator, err := r.config.Storage.GetDKGOperator(msg.Signer)
	if !found {
		return errors.New("unable to find signer")
	}
	if err != nil {
		return errors.Wrap(err, "unable to find signer")
	}

	root, err := msg.Data.GetRoot()
	if err != nil {
		return errors.Wrap(err, "fail to get root")
	}
	pk, err := crypto.Ecrecover(root, msg.Signature)
	if err != nil {
		return errors.New("unable to recover public key")
	}
	addr := common.BytesToAddress(crypto.Keccak256(pk[1:])[12:])
	if addr != operator.ETHAddress {
		return errors.New("invalid signature")
	}
	return nil
}

func (r *Runner) validateDepositDataRoot(msg *PartialDepositData) error {
	if !bytes.Equal(r.DepositDataRoot, msg.Root) {
		return errors.New("deposit data roots not equal")
	}
	return nil
}

func (r *Runner) validateDepositDataSig(msg *PartialDepositData) error {

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

func (r *Runner) ownOutput() *SignedOutput {
	return r.OutputMsgs[r.Operator.OperatorID]
}
