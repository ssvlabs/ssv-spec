package dkg

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type MsgType int

const (
	InitMsgType MsgType = iota
	ProtocolMsgType
	DepositDataMsgType
)

type Message struct {
	MsgType    MsgType
	Identifier types.MessageID
	Data       []byte
}

// Encode returns a msg encoded bytes or error
func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

type SignedMessage struct {
	Message   *Message
	Signer    types.OperatorID
	Signature types.Signature
}

// Encode returns a msg encoded bytes or error
func (signedMsg *SignedMessage) Encode() ([]byte, error) {
	return json.Marshal(signedMsg)
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, signedMsg)
}

// Init is the first message in a DKG which initiates a DKG
type Init struct {
	// OperatorIDs are the operators selected for the DKG
	OperatorIDs []types.OperatorID
	// Threshold DKG threshold for signature reconstruction
	Threshold uint16
	// WithdrawalCredentials used when signing the deposit data
	WithdrawalCredentials []byte
}

// Encode returns a msg encoded bytes or error
func (msg *Init) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Init) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

// Output is the last message in every DKG which marks a specific node's end of process
type Output struct {
	// Identifier of the DKG
	Identifier types.MessageID
	// EncryptedShare standard SSV encrypted shares
	EncryptedShare []byte
	// DKGSize number of participants in the DKG
	DKGSetSize uint16
	// Threshold DKG threshold for signature reconstruction
	Threshold uint16
	// ValidatorPubKey the resulting public key corresponding to the shared private key
	ValidatorPubKey types.ValidatorPK
	// WithdrawalCredentials same as in Init
	WithdrawalCredentials []byte
	// SignedDepositData reconstructed signature of DepositMessage according to eth2 spec
	SignedDepositData types.Signature
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
		o.Threshold,
		o.ValidatorPubKey,
		o.WithdrawalCredentials,
		o.SignedDepositData,
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

// PartialDepositData contains a partial deposit data signature
type PartialDepositData struct {
	Signer    types.OperatorID
	Root      types.Root
	Signature types.Signature
}

// Encode returns a msg encoded bytes or error
func (msg *PartialDepositData) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *PartialDepositData) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
