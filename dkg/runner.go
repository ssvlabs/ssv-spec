package dkg

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type Runner interface {
	ProcessMsg(msg *SignedMessage) (bool, error)
}

// Runner manages the execution of a DKG, start to finish.
type runner struct {
	Operator *Operator
	// InitMsg holds the init method which started this runner
	InitMsg *Init
	// ReshareMsg holds the reshare method which started this runner
	ReshareMsg *Reshare
	// Identifier unique for DKG session
	Identifier RequestID
	// KeygenOutcome holds the protocol outcome once it finishes
	KeygenOutcome *ProtocolOutcome
	// DepositDataRoot is the signing root for the deposit data
	DepositDataRoot []byte
	// DepositDataSignatures holds partial sigs on deposit data
	DepositDataSignatures map[types.OperatorID]*PartialDepositData
	// OutputMsgs holds all output messages received
	OutputMsgs map[types.OperatorID]*SignedOutput

	protocol Protocol
	config   *Config
}

// ProcessMsg processes a DKG signed message and returns true and stream keygen output or blame if finished
func (r *runner) ProcessMsg(msg *SignedMessage) (bool, error) {
	// TODO - validate message

	switch msg.Message.MsgType {
	case ProtocolMsgType:
		if r.DepositDataSignatures[r.Operator.OperatorID] != nil {
			return false, errors.New("keygen has already completed")
		}

		finished, o, err := r.protocol.ProcessMsg(msg)
		if err != nil {
			return false, errors.Wrap(err, "failed to process dkg msg")
		}

		if finished {
			r.KeygenOutcome = o
			isBlame, err := r.KeygenOutcome.IsFailedWithBlame()
			if err != nil {
				return true, errors.Wrap(err, "invalid KeygenOutcome")
			}
			if isBlame {
				err := r.config.Network.StreamDKGBlame(r.KeygenOutcome.BlameOutput)
				return true, errors.Wrap(err, "failed to stream blame output")
			}
			if r.KeygenOutcome.ProtocolOutput == nil {
				return true, errors.Wrap(err, "protocol finished without blame or keygen result")
			}

			if r.isResharing() {
				if err := r.prepareAndBroadcastOutput(); err != nil {
					return false, err
				}
			} else {
				if err := r.prepareAndBroadcastDepositData(); err != nil {
					return false, err
				}
			}

		}
		return false, nil
	case DepositDataMsgType:
		depSig := &PartialDepositData{}
		if err := depSig.Decode(msg.Message.Data); err != nil {
			return false, errors.Wrap(err, "could not decode PartialDepositData")
		}

		if err := r.validateDepositDataSig(depSig); err != nil {
			return false, errors.Wrap(err, "PartialDepositData invalid")
		}

		if found := r.DepositDataSignatures[msg.Signer]; found == nil {
			r.DepositDataSignatures[msg.Signer] = depSig
		} else if !bytes.Equal(found.Signature, msg.Signature) {
			return false, errors.New("inconsistent partial signature received")
		}

		if len(r.DepositDataSignatures) == int(r.InitMsg.Threshold) {
			if err := r.prepareAndBroadcastOutput(); err != nil {
				return false, err
			}
		}
		return false, nil
	case OutputMsgType:
		output := &SignedOutput{}
		if err := output.Decode(msg.Message.Data); err != nil {
			return false, errors.Wrap(err, "could not decode SignedOutput")
		}

		if err := r.validateSignedOutput(output); err != nil {
			return false, errors.Wrap(err, "signed output invali")
		}

		r.OutputMsgs[msg.Signer] = output
		// GLNOTE: Actually we need every operator to sign instead only the quorum!
		finished := false
		if !r.isResharing() {
			finished = len(r.OutputMsgs) == len(r.InitMsg.OperatorIDs)
		} else {
			finished = len(r.OutputMsgs) == len(r.ReshareMsg.OperatorIDs)
		}
		if finished {
			err := r.config.Network.StreamDKGOutput(r.OutputMsgs)
			return true, errors.Wrap(err, "failed to stream dkg output")
		}

		return false, nil
	default:
		return false, errors.New("msg type invalid")
	}
}

func (r *runner) prepareAndBroadcastDepositData() error {
	// generate deposit data
	root, _, err := types.GenerateETHDepositData(
		r.KeygenOutcome.ProtocolOutput.ValidatorPK,
		r.InitMsg.WithdrawalCredentials,
		r.InitMsg.Fork,
		types.DomainDeposit,
	)
	if err != nil {
		return errors.Wrap(err, "could not generate deposit data")
	}

	r.DepositDataRoot = root

	// sign
	sig := r.KeygenOutcome.ProtocolOutput.Share.SignByte(root)

	// broadcast
	pdd := &PartialDepositData{
		Signer:    r.Operator.OperatorID,
		Root:      r.DepositDataRoot,
		Signature: sig.Serialize(),
	}
	if err := r.signAndBroadcastMsg(pdd, DepositDataMsgType); err != nil {
		return errors.Wrap(err, "could not broadcast partial deposit data")
	}
	r.DepositDataSignatures[r.Operator.OperatorID] = pdd
	return nil
}

func (r *runner) prepareAndBroadcastOutput() error {
	var (
		depositSig types.Signature
		err        error
	)
	if r.isResharing() {
		depositSig = nil
	} else {
		// reconstruct deposit data sig
		depositSig, err = r.reconstructDepositDataSignature()
		if err != nil {
			return errors.Wrap(err, "could not reconstruct deposit data sig")
		}
	}

	// encrypt Operator's share
	encryptedShare, err := r.config.Signer.Encrypt(r.Operator.EncryptionPubKey, r.KeygenOutcome.ProtocolOutput.Share.Serialize())
	if err != nil {
		return errors.Wrap(err, "could not encrypt share")
	}

	ret, err := r.generateSignedOutput(&Output{
		RequestID:            r.Identifier,
		EncryptedShare:       encryptedShare,
		SharePubKey:          r.KeygenOutcome.ProtocolOutput.Share.GetPublicKey().Serialize(),
		ValidatorPubKey:      r.KeygenOutcome.ProtocolOutput.ValidatorPK,
		DepositDataSignature: depositSig,
	})
	if err != nil {
		return errors.Wrap(err, "could not generate dkg SignedOutput")
	}

	r.OutputMsgs[r.Operator.OperatorID] = ret
	if err := r.signAndBroadcastMsg(ret, OutputMsgType); err != nil {
		return errors.Wrap(err, "could not broadcast SignedOutput")
	}
	return nil
}

func (r *runner) signAndBroadcastMsg(msg types.Encoder, msgType MsgType) error {
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

func (r *runner) reconstructDepositDataSignature() (types.Signature, error) {
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

func (r *runner) validateSignedOutput(msg *SignedOutput) error {
	// TODO: Separate fields match and signature validation
	output := r.ownOutput()
	if output != nil {
		if output.BlameData == nil {
			if output.Data.RequestID != msg.Data.RequestID {
				return errors.New("got mismatching RequestID")
			}
			if !bytes.Equal(output.Data.ValidatorPubKey, msg.Data.ValidatorPubKey) {
				return errors.New("got mismatching ValidatorPubKey")
			}
		} else {
			if output.BlameData.RequestID != msg.BlameData.RequestID {
				return errors.New("got mismatching RequestID")
			}
		}
	}

	found, operator, err := r.config.Storage.GetDKGOperator(msg.Signer)
	if !found {
		return errors.New("unable to find signer")
	}
	if err != nil {
		return errors.Wrap(err, "unable to find signer")
	}

	var (
		root []byte
	)

	if msg.BlameData == nil {
		root, err = msg.Data.GetRoot()
	} else {
		root, err = msg.BlameData.GetRoot()
	}
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

func (r *runner) validateDepositDataRoot(msg *PartialDepositData) error {
	if !bytes.Equal(r.DepositDataRoot, msg.Root) {
		return errors.New("deposit data roots not equal")
	}
	return nil
}

func (r *runner) validateDepositDataSig(msg *PartialDepositData) error {

	// find operator and verify msg
	sharePK, found := r.KeygenOutcome.ProtocolOutput.OperatorPubKeys[msg.Signer]
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

func (r *runner) generateSignedOutput(o *Output) (*SignedOutput, error) {
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

func (r *runner) ownOutput() *SignedOutput {
	return r.OutputMsgs[r.Operator.OperatorID]
}

func (r *runner) isResharing() bool {
	return r.ReshareMsg != nil
}
