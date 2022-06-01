package dkg

import (
	"crypto/ecdsa"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type MsgType int

const (
	InitMsgType MsgType = iota
	OutputMsgType
)

type Message struct {
	MsgType MsgType
	Data    []byte
}

type SignedMessage struct {
	Message   *Message
	Signer    types.OperatorID
	Signature types.Signature
}

// Init is the first message in a DKG which initiates a DKG
type Init struct {
	// OperatorIDs are the operators selected for the DKG
	OperatorIDs []types.OperatorID
	// WithdrawalCredentials used when signing the deposit data
	WithdrawalCredentials []byte
}

// Output is the last message in every DKG which marks a specific node's end of process
type Output struct {
	// EncryptedShare standard SSV encrypted shares
	EncryptedShare []byte
	// DKGSize number of participants in the DKG
	DKGSetSize uint16
	// ValidatorPubKey the resulting public key corresponding to the shared private key
	ValidatorPubKey types.ValidatorPK
	// WithdrawalCredentials same as in Init
	WithdrawalCredentials []byte
	// PartialSignedDepositData partial signature of DepositMessage according to eth2 spec
	PartialSignedDepositData types.Signature
}

func (o *Output) GetRoot() ([]byte, error) {
	uint16Solidity, _ := abi.NewType("uint16", "", nil)
	bytesSolidity, _ := abi.NewType("bytes", "", nil)

	arguments := abi.Arguments{
		{
			Type: bytesSolidity,
		},
		{
			Type: uint16Solidity,
		},
		{
			Type: bytesSolidity,
		},
		{
			Type: bytesSolidity,
		},
		{
			Type: bytesSolidity,
		},
	}

	bytes, _ := arguments.Pack(
		o.EncryptedShare,
		o.DKGSetSize,
		o.ValidatorPubKey,
		o.WithdrawalCredentials,
		o.PartialSignedDepositData,
	)

	return crypto.Keccak256(bytes), nil
}

type SignedOutput struct {
	// Data signed
	Data *Output
	// Signer operator ID which signed
	Signer types.OperatorID
	// Signature over Data.GetRoot()
	Signature types.Signature
}

func SignOutput(output *Output, privKey *ecdsa.PrivateKey) (types.Signature, error) {
	root, err := output.GetRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not get root from output message")
	}

	return crypto.Sign(root, privKey)
}
