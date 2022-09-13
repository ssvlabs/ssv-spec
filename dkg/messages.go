package dkg

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type RequestID [24]byte

const (
	ethAddressSize     = 20
	ethAddressStartPos = 0
	indexSize          = 4
	indexStartPos      = ethAddressStartPos + ethAddressSize
)

func (msg RequestID) GetETHAddress() common.Address {
	ret := common.Address{}
	copy(ret[:], msg[ethAddressStartPos:ethAddressStartPos+ethAddressSize])
	return ret
}

func (msg RequestID) GetRoleType() uint32 {
	indexByts := msg[indexStartPos : indexStartPos+indexSize]
	return binary.LittleEndian.Uint32(indexByts)
}

func NewRequestID(ethAddress common.Address, index uint32) RequestID {
	indexByts := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexByts, index)

	ret := RequestID{}
	copy(ret[ethAddressStartPos:ethAddressStartPos+ethAddressSize], ethAddress[:])
	copy(ret[indexStartPos:indexStartPos+indexSize], indexByts[:])
	return ret
}

type MsgType int

const (
	// InitMsgType sent when DKG instance is started by requester
	InitMsgType MsgType = iota
	// ProtocolMsgType is the DKG itself
	ProtocolMsgType
	// DepositDataMsgType post DKG deposit data signatures
	DepositDataMsgType
	// OutputMsgType final output msg used by requester to make deposits and register validator with SSV
	OutputMsgType
)

type Message struct {
	MsgType    MsgType
	Identifier RequestID
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

func (msg *Message) Validate() error {
	// TODO msg type
	// TODO len(data)
	return nil
}

func (msg *Message) GetRoot() ([]byte, error) {
	marshaledRoot, err := msg.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode PartialSignatureMessage")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
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

func (signedMsg *SignedMessage) Validate() error {
	// TODO len(sig) == ecdsa sig lenth

	return signedMsg.Message.Validate()
}

func (signedMsg *SignedMessage) GetRoot() ([]byte, error) {
	return signedMsg.Message.GetRoot()
}

// Init is the first message in a DKG which initiates a DKG
type Init struct {
	// OperatorIDs are the operators selected for the DKG
	OperatorIDs []types.OperatorID
	// Threshold DKG threshold for signature reconstruction
	Threshold uint16
	// WithdrawalCredentials used when signing the deposit data
	WithdrawalCredentials []byte
	// Fork is eth2 fork version
	Fork phase0.Version
}

func (msg *Init) Validate() error {
	if len(msg.WithdrawalCredentials) != phase0.HashLength {
		return errors.New("invalid WithdrawalCredentials")
	}
	if len(msg.OperatorIDs) < 4 || (len(msg.OperatorIDs)-1)%3 != 0 {
		return errors.New("invalid number of operators which has to be 3f+1")
	}

	if int(msg.Threshold) != (len(msg.OperatorIDs)-1)*2/3+1 {
		return errors.New("invalid threshold which has to be 2f+1")
	}

	return nil
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
	// RequestID for the DKG instance (not used for signing)
	RequestID RequestID
	// EncryptedShare standard SSV encrypted shares
	EncryptedShare []byte
	// SharePubKey is the share's BLS pubkey
	SharePubKey []byte
	// ValidatorPubKey the resulting public key corresponding to the shared private key
	ValidatorPubKey types.ValidatorPK
	// DepositDataSignature reconstructed signature of DepositMessage according to eth2 spec
	DepositDataSignature types.Signature
}

func (o *Output) GetRoot() ([]byte, error) {
	bytesSolidity, _ := abi.NewType("bytes", "", nil)

	arguments := abi.Arguments{
		{
			Type: bytesSolidity,
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

	bytes, err := arguments.Pack(
		[]byte(o.EncryptedShare),
		[]byte(o.SharePubKey),
		[]byte(o.ValidatorPubKey),
		[]byte(o.DepositDataSignature),
	)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(bytes), nil
}

type SignedOutput struct {
	// Data signed
	Data *Output
	// Signer Operator ID which signed
	Signer types.OperatorID
	// Signature over Data.GetRoot()
	Signature types.Signature
}

// Encode returns a msg encoded bytes or error
func (msg *SignedOutput) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *SignedOutput) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
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
	Root      []byte
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
