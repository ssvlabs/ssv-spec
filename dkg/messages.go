package dkg

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
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
	InitMsgType MsgType = iota
	ProtocolMsgType
	KeygenOutputType
	PartialSigType
	DepositDataMsgType
	PartialOutputMsgType
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
	// Nonce is used to differentiate DKG tasks of the same OperatorIDs and WithdrawalCredentials
	Nonce int64
	// OperatorIDs are the operators selected for the DKG
	OperatorIDs []types.OperatorID
	// Threshold DKG threshold for signature reconstruction
	Threshold uint16
	// WithdrawalCredentials used when signing the deposit data
	WithdrawalCredentials []byte
}

func (msg *Init) Validate() error {
	// TODO len(operators == 4,7,10,13
	// threshold equal to 2/3 of 4,7,10,13
	// len(WithdrawalCredentials) is valid
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

// TODO: What's the difference / intention of this vs Output.
type KeygenOutput struct {
	Index           uint16
	Threshold       uint16
	ShareCount      uint16
	PublicKey       []byte
	SecretShare     []byte
	SharePublicKeys [][]byte
}

// Encode returns a msg encoded bytes or error
func (d *KeygenOutput) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *KeygenOutput) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

type PartialSignature struct {
	I      uint16
	SigmaI spec.BLSSignature
}

// Encode returns a msg encoded bytes or error
func (d *PartialSignature) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *PartialSignature) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Output is the last message in every DKG which marks a specific node's end of process
type Output struct {
	// Identifier of the DKG
	Identifier RequestID
	// ShareIndex the 1-based index of the share
	ShareIndex uint16
	// EncryptedShare standard SSV encrypted shares
	EncryptedShare []byte
	// SharePubKeys the public keys corresponding to the shares
	SharePubKeys [][]byte
	// DKGSetSize number of participants in the DKG
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

	// TODO: Include RequestID, SharePubKeys and ShareIndex
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
	// Signer Operator ID which signed
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
