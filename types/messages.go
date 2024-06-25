package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ValidatorPK is an eth2 validator public key 48 bytes long
type ValidatorPK phase0.BLSPubKey

// ShareValidatorPK is a partial eth2 validator public key 48 bytes long
type ShareValidatorPK []byte

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
	toMatch := msgID.GetDutyExecutorID()
	return bytes.Equal(vid[:], toMatch)
}

// MessageID is used to identify and route messages to the right validator and Runner
type MessageID [56]byte

func (msg MessageID) GetDomain() []byte {
	return msg[domainStartPos : domainStartPos+domainSize]
}

func (msg MessageID) GetDutyExecutorID() []byte {
	return msg[dutyExecutorIDStartPos : dutyExecutorIDStartPos+dutyExecutorIDSize]
}

func (msg MessageID) GetRoleType() RunnerRole {
	roleByts := msg[roleTypeStartPos : roleTypeStartPos+roleTypeSize]
	return RunnerRole(binary.LittleEndian.Uint32(roleByts))
}

func NewMsgID(domain DomainType, pk []byte, role RunnerRole) MessageID {
	roleByts := make([]byte, 4)
	binary.LittleEndian.PutUint32(roleByts, uint32(role))

	return newMessageID(domain[:], roleByts, pk)
}

func (msgID MessageID) String() string {
	return hex.EncodeToString(msgID[:])
}

func newMessageID(domain, roleByts, dutyExecutorID []byte) MessageID {
	mid := MessageID{}
	copy(mid[domainStartPos:domainStartPos+domainSize], domain[:])
	copy(mid[roleTypeStartPos:roleTypeStartPos+roleTypeSize], roleByts)
	prefixLen := dutyExecutorIDSize - len(dutyExecutorID)
	copy(mid[dutyExecutorIDStartPos+prefixLen:dutyExecutorIDStartPos+dutyExecutorIDSize], dutyExecutorID)
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
	// Data max size is the max between max(qbft.SignedMessage) and max(PartialSignatureMessages)
	// i.e., = max(722412, 144020) = 722412
	Data []byte `ssz-max:"722412"`
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

// SignedSSVMessage is the main message passed within the SSV network. It encapsulates the SSVMessage structure and a signature
type SignedSSVMessage struct {
	Signature  [256]byte // Created by the operator's network key
	OperatorID OperatorID
	Data       []byte
}

// GetOperatorID returns the sender operator ID
func (msg *SignedSSVMessage) GetOperatorID() OperatorID {
	return msg.OperatorID
}

// GetSignature returns the signature of the OperatorID over Data
func (msg *SignedSSVMessage) GetSignature() [256]byte {
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
	return EncodeSignedSSVMessage(msg.Data, msg.OperatorID, msg.Signature), nil
}

// Decode returns error if decoding failed
func (msg *SignedSSVMessage) Decode(data []byte) error {
	msgData, operatorID, signature, err := DecodeSignedSSVMessage(data)
	if err != nil {
		return errors.Wrap(err, "could not decode data into a SignedSSVMessage")
	}
	msg.Data = msgData
	msg.OperatorID = operatorID
	msg.Signature = signature
	return nil
}

// Validate checks the following rules:
// - OperatorID should not be 0
// - Signature length should not be 0
// - Data length should not be 0
func (msg *SignedSSVMessage) Validate() error {
	if msg.OperatorID == 0 {
		return errors.New("signer ID 0 not allowed")
	}
	if len(msg.Data) == 0 {
		return errors.New("Data has length 0 in SignedSSVMessage")
	}
	return nil
}

// EncodeSignedSSVMessage serializes the message, op id and signature into bytes
func EncodeSignedSSVMessage(message []byte, operatorID OperatorID, signature [256]byte) []byte {
	b := make([]byte, signatureSize+operatorIDSize+len(message))
	copy(b[signatureOffset:], signature[:])
	binary.LittleEndian.PutUint64(b[operatorIDOffset:], operatorID)
	copy(b[messageOffset:], message)
	return b
}

// DecodeSignedSSVMessage deserializes signed message bytes messsage, op id and a signature
func DecodeSignedSSVMessage(encoded []byte) ([]byte, OperatorID, [256]byte, error) {
	if len(encoded) < messageOffset {
		return nil, 0, [256]byte{}, fmt.Errorf("unexpected encoded message size of %d", len(encoded))
	}

	message := encoded[messageOffset:]
	operatorID := binary.LittleEndian.Uint64(encoded[operatorIDOffset : operatorIDOffset+operatorIDSize])
	signature := [256]byte{}
	copy(signature[:], encoded[signatureOffset:signatureOffset+signatureSize])
	return message, operatorID, signature, nil
}

func SSVMessageToSignedSSVMessage(msg *SSVMessage, operatorID OperatorID, signSSVMessageF SignSSVMessageF) (*SignedSSVMessage, error) {
	encodedSSVMsg, err := msg.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode SSVMessage")
	}

	sig, err := signSSVMessageF(encodedSSVMsg)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign SSVMessage")
	}

	return &SignedSSVMessage{
		Signature:  sig,
		OperatorID: operatorID,
		Data:       encodedSSVMsg,
	}, nil
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (signedMsg *SignedSSVMessage) DeepCopy() *SignedSSVMessage {
	ret := &SignedSSVMessage{
		OperatorIDs: make([]OperatorID, len(signedMsg.GetOperatorIDs())),
		Signatures:  make([][]byte, len(signedMsg.Signatures)),
	}
	copy(ret.OperatorIDs, signedMsg.GetOperatorIDs())
	copy(ret.Signatures, signedMsg.Signatures)

	retSSV := &SSVMessage{
		MsgType: signedMsg.SSVMessage.MsgType,
		Data:    make([]byte, len(signedMsg.SSVMessage.Data)),
	}
	msgID := [56]byte{}
	copy(msgID[:], signedMsg.SSVMessage.MsgID[:])
	retSSV.MsgID = msgID

	copy(retSSV.Data, signedMsg.SSVMessage.Data)

	if len(signedMsg.FullData) > 0 {
		ret.FullData = make([]byte, len(signedMsg.FullData))
		copy(ret.FullData, signedMsg.FullData)
	}
	ret.SSVMessage = retSSV
	return ret
}

// MatchedSigners returns true if the provided signer ids are equal to GetOperatorIDs() without order significance
func (msg *SignedSSVMessage) MatchedSigners(ids []OperatorID) bool {
	if len(msg.GetOperatorIDs()) != len(ids) {
		return false
	}

	for _, id := range msg.GetOperatorIDs() {
		found := false
		for _, id2 := range ids {
			if id == id2 {
				found = true
			}
		}

		if !found {
			return false
		}
	}
	return true
}

// CommonSigners returns true if there is at least 1 common signer
func (msg *SignedSSVMessage) CommonSigners(ids []OperatorID) bool {
	for _, id := range msg.GetOperatorIDs() {
		for _, id2 := range ids {
			if id == id2 {
				return true
			}
		}
	}
	return false
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (msg *SignedSSVMessage) Aggregate(msgToAggregate *SignedSSVMessage) error {
	if msg.CommonSigners(msgToAggregate.GetOperatorIDs()) {
		return errors.New("duplicate signers")
	}

	e1, err := msg.SSVMessage.Encode()
	if err != nil {
		return errors.Wrap(err, "could not get own encoded SSVMessage")
	}
	e2, err := msgToAggregate.SSVMessage.Encode()
	if err != nil {
		return errors.Wrap(err, "could not get encoded SSVMessage to be aggregated")
	}
	if !bytes.Equal(e1[:], e2[:]) {
		return errors.New("can't aggregate, encoded messages not equal")
	}

	msg.Signatures = append(msg.Signatures, msgToAggregate.Signatures...)
	msg.OperatorIDs = append(msg.OperatorIDs, msgToAggregate.GetOperatorIDs()...)

	return nil
}

// Check if all signedMsg's signers belong to the given committee in O(n+m)
func (msg *SignedSSVMessage) CheckSignersInCommittee(operators []*Operator) bool {
	// Committee's operators map
	committeeMap := make(map[OperatorID]struct{})
	for _, op := range operators {
		committeeMap[op.OperatorID] = struct{}{}
	}

	// Check that all message signers belong to the map
	for _, signer := range msg.GetOperatorIDs() {
		if _, ok := committeeMap[signer]; !ok {
			return false
		}
	}
	return true
}

// WithoutFullData returns SignedMessage without full data
func (msg *SignedSSVMessage) WithoutFullData() *SignedSSVMessage {
	return &SignedSSVMessage{
		OperatorIDs: msg.OperatorIDs,
		Signatures:  msg.Signatures,
		SSVMessage:  msg.SSVMessage,
	}
}
