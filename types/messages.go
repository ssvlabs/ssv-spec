package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"
)

// ValidatorPK is an eth2 validator public key
type ValidatorPK []byte

const (
	domainSize       = 4
	domainStartPos   = 0
	pubKeySize       = 48
	pubKeyStartPos   = domainStartPos + domainSize
	roleTypeSize     = 4
	roleTypeStartPos = pubKeyStartPos + pubKeySize
)

type Validate interface {
	// Validate returns error if msg validation doesn't pass.
	// Msg validation checks the msg, it's variables for validity.
	Validate() error
}

// MessageIDBelongs returns true if message ID belongs to validator
func (vid ValidatorPK) MessageIDBelongs(msgID MessageID) bool {
	toMatch := msgID.GetPubKey()
	return bytes.Equal(vid, toMatch)
}

// MessageID is used to identify and route messages to the right validator and Runner
type MessageID [56]byte

func (msg MessageID) GetDomain() []byte {
	return msg[domainStartPos : domainStartPos+domainSize]
}

func (msg MessageID) GetPubKey() []byte {
	return msg[pubKeyStartPos : pubKeyStartPos+pubKeySize]
}

func (msg MessageID) GetRoleType() BeaconRole {
	roleByts := msg[roleTypeStartPos : roleTypeStartPos+roleTypeSize]
	return BeaconRole(binary.LittleEndian.Uint32(roleByts))
}

func NewMsgID(domain DomainType, pk []byte, role BeaconRole) MessageID {
	roleByts := make([]byte, 4)
	binary.LittleEndian.PutUint32(roleByts, uint32(role))

	return newMessageID(domain[:], pk, roleByts)
}

func (msgID MessageID) String() string {
	return hex.EncodeToString(msgID[:])
}

func MessageIDFromBytes(mid []byte) MessageID {
	if len(mid) < domainSize+pubKeySize+roleTypeSize {
		return MessageID{}
	}
	return newMessageID(
		mid[domainStartPos:domainStartPos+domainSize],
		mid[pubKeyStartPos:pubKeyStartPos+pubKeySize],
		mid[roleTypeStartPos:roleTypeStartPos+roleTypeSize],
	)
}

func newMessageID(domain, pk, roleByts []byte) MessageID {
	mid := MessageID{}
	copy(mid[domainStartPos:domainStartPos+domainSize], domain[:])
	copy(mid[pubKeyStartPos:pubKeyStartPos+pubKeySize], pk)
	copy(mid[roleTypeStartPos:roleTypeStartPos+roleTypeSize], roleByts)
	return mid
}

type MsgType uint64

const (
	// SSVConsensusMsgType are all QBFT consensus related messages
	SSVConsensusMsgType MsgType = iota
	// SSVPartialSignatureMsgType are all partial signatures msgs over beacon chain specific signatures
	SSVPartialSignatureMsgType
	// DKGMsgType represent all DKG related messages
	DKGMsgType
)

// MessageSignature includes all functions relevant for a signed message (QBFT message, post consensus msg, etc)
type MessageSignature interface {
	Root
	GetSignature() Signature
	GetSigners() []OperatorID
}

// SSVMessage is the main message passed within the SSV network, it can contain different types of messages (QBTF, Sync, etc.)
type SSVMessage struct {
	MsgType MsgType
	MsgID   MessageID `ssz-size:"56"`
	// Data max size is qbft SignedMessage max ~= 2^22 + 2^20 + 96 + 13 + 2^20 ~= 2^23
	Data []byte `ssz-max:"8388608"` // 2^23
}

func (msg *SSVMessage) GetType() MsgType {
	return msg.MsgType
}

// GetID returns a unique msg ID that is used to identify to which validator should the message be sent for processing
func (msg *SSVMessage) GetID() MessageID {
	return msg.MsgID
}

// GetData returns message Data as byte slice
func (msg *SSVMessage) GetData() []byte {
	return msg.Data
}

// Encode returns a msg encoded bytes or error
func (msg *SSVMessage) Encode() ([]byte, error) {
	return msg.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msg *SSVMessage) Decode(data []byte) error {
	return msg.UnmarshalSSZ(data)
}

// SSVMessage is the main message passed within the SSV network. It encapsulates the SSVMessage structure and a signature
type SignedSSVMessage struct {
	OperatorID OperatorID
	Signature  []byte `ssz-max:"512"`     // Current signature max size allow keys up to 512*8 = 4096 bits
	Data       []byte `ssz-max:"8388672"` // Max size extracted from SSVMessage
}

// GetOperatorID returns the sender operator ID
func (msg *SignedSSVMessage) GetOperatorID() OperatorID {
	return msg.OperatorID
}

// GetSignature returns the signature of the OperatorID over Data
func (msg *SignedSSVMessage) GetSignature() []byte {
	return msg.Signature
}

// GetData returns message Data as byte slice
func (msg *SignedSSVMessage) GetData() []byte {
	return msg.Data
}

// GetSSVMessageFromData returns message SSVMessage decoded from data
func (msg *SignedSSVMessage) GetSSVMessageFromData() (*SSVMessage, error) {
	ssvMessage := &SSVMessage{}
	err := ssvMessage.Decode(msg.Data)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode SSVMessage from data in SignedSSVMessage")
	}
	return ssvMessage, nil
}

// Encode returns a msg encoded bytes or error
func (msg *SignedSSVMessage) Encode() ([]byte, error) {
	return msg.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msg *SignedSSVMessage) Decode(data []byte) error {
	return msg.UnmarshalSSZ(data)
}

// Validate checks the following rules:
// - OperatorID should not be 0
// - Signature length should not be 0
// - Data length should not be 0
func (msg *SignedSSVMessage) Validate() error {
	if msg.OperatorID == 0 {
		return errors.New("OperatorID in SignedSSVMessage is 0")
	}
	if len(msg.Signature) == 0 {
		return errors.New("Signature has length 0 in SignedSSVMessage")
	}
	if len(msg.Data) == 0 {
		return errors.New("Data has length 0 in SignedSSVMessage")
	}
	return nil
}
