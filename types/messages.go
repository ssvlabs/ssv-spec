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
	// Data max size is qbft SignedMessage max ~= 5243144 + 2^20 + 96 + 13 ~= 6291829
	Data []byte `ssz-max:"6291829"`
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
	OperatorID []OperatorID `ssz-max:"13"`
	Signature  [][]byte     `ssz-max:"13,512"`  // Created by the operator's private key. Max size allow keys up to 512*8 = 4096 bits
	SSVMessage *SSVMessage  `ssz-max:"6291893"` // Max size extracted from SSVMessage
}

// GetOperatorID returns the sender operator ID
func (msg *SignedSSVMessage) GetOperatorIDs() []OperatorID {
	return msg.OperatorID
}

// GetSignature returns the signature of the OperatorID over Data
func (msg *SignedSSVMessage) GetSignatures() [][]byte {
	return msg.Signature
}

// GetData returns message Data as byte slice
func (msg *SignedSSVMessage) GetSSVMessage() *SSVMessage {
	return msg.SSVMessage
}

// Encode returns a msg encoded bytes or error
func (msg *SignedSSVMessage) Encode() ([]byte, error) {
	return msg.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msg *SignedSSVMessage) Decode(data []byte) error {
	return msg.UnmarshalSSZ(data)
}

// GetRoot returns the root
func (msg *SignedSSVMessage) GetRoot() ([32]byte, error) {
	return msg.HashTreeRoot()
}

// Validate
func (msg *SignedSSVMessage) Validate() error {
	// There must be at least one signer
	if len(msg.OperatorID) == 0 {
		return errors.New("No OperatorID in SignedSSVMessage")
	}
	// Each signer must be different than 0 and unique
	operatorsSet := make(map[OperatorID]bool)
	for _, operatorID := range msg.OperatorID {
		if operatorID == 0 {
			return errors.New("OperatorID in SignedSSVMessage is 0")
		}
		if operatorsSet[operatorID] {
			return errors.New("non unique signer")
		}
		operatorsSet[operatorID] = true
	}
	// There must be at least one signature
	if len(msg.Signature) == 0 {
		return errors.New("No signature in SignedSSVMessage")
	}
	// No signature can be empty
	for _, signature := range msg.Signature {
		if len(signature) == 0 {
			return errors.New("Signature has length 0 in SignedSSVMessage")
		}
	}
	// There must be an equal number of signers and signatures
	if len(msg.OperatorID) != len(msg.Signature) {
		return errors.New("SignedSSVMessage has a different number of operato IDs and signatures")
	}
	// The SSVMessage can't be nil
	if msg.SSVMessage == nil {
		return errors.New("SSVMessage is nil")
	}
	return nil
}
