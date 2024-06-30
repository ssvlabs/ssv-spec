package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ValidatorPK is an eth2 validator public key 48 bytes long
type ValidatorPK phase0.BLSPubKey

// ShareValidatorPK is a partial eth2 validator public key 48 bytes long
type ShareValidatorPK []byte

const (
	domainSize             = 4
	domainStartPos         = 0
	roleTypeSize           = 4
	roleTypeStartPos       = domainStartPos + domainSize
	dutyExecutorIDSize     = 48
	dutyExecutorIDStartPos = roleTypeStartPos + roleTypeSize
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

func NewMsgID(domain DomainType, dutyExecutorID []byte, role RunnerRole) MessageID {
	roleByts := make([]byte, 4)
	binary.LittleEndian.PutUint32(roleByts, uint32(role))

	return newMessageID(domain[:], roleByts, dutyExecutorID)
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
	Signatures  [][]byte     `ssz-max:"13,256"` // Created by the operators' key
	OperatorIDs []OperatorID `ssz-max:"13"`
	SSVMessage  *SSVMessage
	// Full data max value is the max value between ValidatorConsensusData and BeaconVote
	FullData []byte `ssz-max:"4194532"`
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

// Validate checks the following rules:
// - OperatorID must have at least one element
// - Any OperatorID must not be 0
// - The number of signatures and OperatorIDs must be the same
// - Any signature must not have length 0
// - SSVMessage must not be nil
func (msg *SignedSSVMessage) Validate() error {
	// Validate OperatorID field
	if len(msg.OperatorIDs) == 0 {
		return errors.New("no signers")
	}
	// Check unique signers
	signed := make(map[OperatorID]struct{})
	for _, operatorID := range msg.OperatorIDs {
		if _, exists := signed[operatorID]; exists {
			return errors.New("non unique signer")
		}
		if operatorID == 0 {
			return errors.New("signer ID 0 not allowed")
		}

		signed[operatorID] = struct{}{}
	}
	// Validate Signature field
	if len(msg.Signatures) == 0 {
		return errors.New("no signatures")
	}
	for _, signature := range msg.Signatures {
		if len(signature) == 0 {
			return errors.New("empty signature")
		}
	}
	// Check that the numbers of signatures and signers are equal
	if len(msg.Signatures) != len(msg.OperatorIDs) {
		return errors.New("number of signatures is different than number of signers")
	}
	// Validate SSVMessage
	if msg.SSVMessage == nil {
		return errors.New("nil SSVMessage")
	}

	return nil
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (signedMsg *SignedSSVMessage) DeepCopy() *SignedSSVMessage {
	ret := &SignedSSVMessage{
		OperatorIDs: make([]OperatorID, len(signedMsg.OperatorIDs)),
		Signatures:  make([][]byte, len(signedMsg.Signatures)),
	}
	copy(ret.OperatorIDs, signedMsg.OperatorIDs)
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
	if len(msg.OperatorIDs) != len(ids) {
		return false
	}

	for _, id := range msg.OperatorIDs {
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
	for _, id := range msg.OperatorIDs {
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
	if msg.CommonSigners(msgToAggregate.OperatorIDs) {
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
	msg.OperatorIDs = append(msg.OperatorIDs, msgToAggregate.OperatorIDs...)

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
	for _, signer := range msg.OperatorIDs {
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
