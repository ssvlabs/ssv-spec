package dkg

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type Config struct {
	// Protocol the DKG protocol implementation
	Protocol Protocol
	Network  Network
	// OperatorID the node's operator ID
	OperatorID types.OperatorID
	// Identifier unique for DKG session
	Identifier types.MessageID
	Signer     types.DKGSigner
	// PubKey signing key for all message
	PubKey *ecdsa.PublicKey
	// EncryptionPubKey encryption pubkey for shares
	EncryptionPubKey *rsa.PublicKey
}

// Runner manages the execution of a DKG, start to finish.
type Runner struct {
	Operators             []types.OperatorID
	Threshold             uint16
	WithdrawalCredentials []byte

	ProtocolOutput        *ProtocolOutput
	DepositDataSignatures map[types.OperatorID]*PartialDepositData

	protocol Protocol
	config   *Config
}

func NewRunner(initMsg *Init, config *Config) (*Runner, error) {
	runner := &Runner{
		Operators:             initMsg.OperatorIDs,
		Threshold:             initMsg.Threshold,
		WithdrawalCredentials: initMsg.WithdrawalCredentials,

		DepositDataSignatures: map[types.OperatorID]*PartialDepositData{},

		protocol: config.Protocol,
		config:   config,
	}

	if err := runner.protocol.Start(initMsg); err != nil {
		return nil, errors.Wrap(err, "could not start dkg protocol")
	}

	return runner, nil
}

// ProcessMsg processes a DKG signed message and returns true and signed output if finished
func (r *Runner) ProcessMsg(msg *SignedMessage) (bool, *SignedOutput, error) {
	// TODO - validate message

	switch msg.Message.MsgType {
	case ProtocolMsgType:
		finished, o, err := r.protocol.ProcessMsg(msg)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to process dkg msg")
		}

		if finished {
			r.ProtocolOutput = o
		}

		// TODO broadcast partial deposit data
	case DepositDataMsgType:
		// TODO validate (including which operator it is)

		depSig := &PartialDepositData{}
		if err := depSig.Decode(msg.Message.Data); err != nil {
			return false, nil, errors.Wrap(err, "could not decode PartialDepositData")
		}

		r.DepositDataSignatures[msg.Signer] = depSig
		if len(r.DepositDataSignatures) >= int(r.Threshold) {
			// reconstruct deposit data sig
			depositSig, err := r.reconstructDepositDataSignature()
			if err != nil {
				return false, nil, errors.Wrap(err, "could not reconstruct deposit data sig")
			}

			// encrypt operator's share
			encryptedShare, err := r.config.Signer.Encrypt(r.config.EncryptionPubKey, r.ProtocolOutput.Share.Serialize())
			if err != nil {
				return false, nil, errors.Wrap(err, "could not encrypt share")
			}

			ret, err := r.generateSignedOutput(&Output{
				Identifier:            r.config.Identifier,
				EncryptedShare:        encryptedShare,
				DKGSetSize:            uint16(len(r.Operators)),
				ValidatorPubKey:       r.ProtocolOutput.ValidatorPK,
				WithdrawalCredentials: r.WithdrawalCredentials,
				SignedDepositData:     depositSig,
			})
			if err != nil {
				return false, nil, errors.Wrap(err, "could not generate dkg SignedOutput")
			}
			return true, ret, nil
		}
	default:
		return false, nil, errors.New("msg type invalid")
	}

	return false, nil, nil
}

func (r *Runner) reconstructDepositDataSignature() (types.Signature, error) {
	panic("implement")
}

func (r *Runner) generateSignedOutput(o *Output) (*SignedOutput, error) {
	sig, err := r.config.Signer.SignDKGOutput(o, r.config.PubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign output")
	}

	return &SignedOutput{
		Data:      o,
		Signer:    r.config.OperatorID,
		Signature: sig,
	}, nil
}
