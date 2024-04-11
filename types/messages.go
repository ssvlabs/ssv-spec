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

	// SignedSSVMessage offsets
	signatureSize    = 256
	signatureOffset  = 0
	operatorIDSize   = 8
	operatorIDOffset = signatureOffset + signatureSize
	messageOffset    = operatorIDOffset + operatorIDSize
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
	Signature  [][]byte     `ssz-max:"13,256"` // Created by the operators' key
	OperatorID []OperatorID `ssz-max:"13"`
	SSVMessage *SSVMessage
	FullData   []byte `ssz-max:"6291829"`
}

// GetOperatorID returns the sender operator ID
func (msg *SignedSSVMessage) GetOperatorID() []OperatorID {
	return msg.OperatorID
}

// GetSignature returns the signature of the OperatorID over Data
func (msg *SignedSSVMessage) GetSignature() [][]byte {
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

// Validate checks the following rules:
// - OperatorID must have at least one element
// - Any OperatorID must not be 0
// - The number of signatures and OperatorIDs must be the same
// - Any signature must not have length 0
// - SSVMessage must not be nil
func (msg *SignedSSVMessage) Validate() error {
	// Validate OperatorID field
	if len(msg.OperatorID) == 0 {
		return errors.New("no signers")
	}
	for _, operatorID := range msg.OperatorID {
		if operatorID == 0 {
			return errors.New("signer ID 0 not allowed")
		}
	}
	// Validate Signature field
	if len(msg.Signature) == 0 {
		return errors.New("no signatures")
	}
	for _, signature := range msg.Signature {
		if len(signature) == 0 {
			return errors.New("empty signature")
		}
	}
	// Check that the numbers of signatures and signers are equal
	if len(msg.Signature) != len(msg.OperatorID) {
		return errors.New("number of signatures is different than number of signers")
	}
	// Validate SSVMessage
	if msg.SSVMessage == nil {
		return errors.New("nil SSVMessage")
	}
	return nil
}

func SSVMessageToSignedSSVMessage(msg *SSVMessage, operatorID OperatorID, signSSVMessageF SignSSVMessageF) (*SignedSSVMessage, error) {

	sig, err := signSSVMessageF(msg)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign SSVMessage")
	}

	return &SignedSSVMessage{
		Signature:  [][]byte{sig},
		OperatorID: []OperatorID{operatorID},
		SSVMessage: msg,
	}, nil
}
