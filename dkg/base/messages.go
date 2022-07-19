package base

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/golang/protobuf/proto"
)

type MsgType int


type SessionId = []byte

type Signable interface {
	types.Root
	SetSignature([]byte) error
}

func toRequestID(id SessionId) RequestID {
	// TODO: Check size
	reqID := new(RequestID)
	copy(reqID[:], id)
	return *reqID
}

// Encode returns a msg encoded bytes or error
func (x *MessageHeader) RequestID() RequestID {
	return toRequestID(x.SessionId)
}

// Encode returns a msg encoded bytes or error
func (x *Message) Encode() ([]byte, error) {
	return proto.Marshal(x)
}

// Decode returns error if decoding failed
func (x *Message) Decode(data []byte) error {
	return proto.Unmarshal(data, x)
}

func (x *Message) Validate() error {
	// TODO: Implement
	return nil
}

func (x *Message) GetRoot() ([]byte, error) {
	raw, err := x.Encode()
	if err != nil {
		return nil, err
	}
	newMsg := &Message{}
	err = newMsg.Decode(raw)
	if err != nil {
		return nil, err
	}
	newMsg.Signature = nil
	bytes, err := newMsg.Encode()
	if err != nil {
		return nil, err
	}
	var root []byte
	rootFixed := sha256.Sum256(bytes)
	copy(root, rootFixed[:])

	return root, nil
}

func (x *Message) SetSignature(bytes []byte) error {
	x.Signature = bytes
	return nil
}

// SignedMessage Deprecated
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
