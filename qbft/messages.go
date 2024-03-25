package qbft

import (
	"crypto/sha256"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// HashDataRoot hashes input data to root
func HashDataRoot(data []byte) ([32]byte, error) {
	ret := sha256.Sum256(data)
	return ret, nil
}

// HasQuorum returns true if a unique set of signers has quorum
func HasQuorum(share *types.Share, msgs []*types.SignedSSVMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.GetOperatorIDs() {
			uniqueSigners[signer] = true
		}
	}
	return share.HasQuorum(len(uniqueSigners))
}

// HasPartialQuorum returns true if a unique set of signers has partial quorum
func HasPartialQuorum(share *types.Share, msgs []*types.SignedSSVMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.GetOperatorIDs() {
			uniqueSigners[signer] = true
		}
	}
	return share.HasPartialQuorum(len(uniqueSigners))
}

type MessageType uint64

const (
	ProposalMsgType MessageType = iota
	PrepareMsgType
	CommitMsgType
	RoundChangeMsgType
)

type Message struct {
	MsgType    MessageType
	Height     Height // QBFT instance Height
	Round      Round  // QBFT round for which the msg is for
	Identifier []byte `ssz-max:"56"` // instance Identifier this msg belongs to

	Root                     [32]byte `ssz-size:"32"`
	DataRound                Round    // The last round that obtained a Prepare quorum
	RoundChangeJustification [][]byte `ssz-max:"13,65536"` // 2^16
	PrepareJustification     [][]byte `ssz-max:"13,65536"` // 2^16

	// Full data max value is ConsensusData max value ~= 2^8 + 8 + 2^20 + 2^22 = 5243144
	FullData []byte `ssz-max:"5243144"`
}

func (msg *Message) GetRoundChangeJustifications() ([]*types.SignedSSVMessage, error) {
	return unmarshalJustifications(msg.RoundChangeJustification)
}

func (msg *Message) GetPrepareJustifications() ([]*types.SignedSSVMessage, error) {
	return unmarshalJustifications(msg.PrepareJustification)
}

func unmarshalJustifications(data [][]byte) ([]*types.SignedSSVMessage, error) {
	ret := make([]*types.SignedSSVMessage, len(data))
	for i, d := range data {
		sMsg := &types.SignedSSVMessage{}
		if err := sMsg.UnmarshalSSZ(d); err != nil {
			return nil, err
		}
		ret[i] = sMsg
	}
	return ret, nil
}

func MarshalJustifications(msgs []*types.SignedSSVMessage) ([][]byte, error) {
	ret := make([][]byte, len(msgs))
	for i, m := range msgs {
		msgWithoutFullData, err := SignedSSVMessageWithoutFullData(m)
		if err != nil {
			return nil, err
		}
		d, err := msgWithoutFullData.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		ret[i] = d
	}
	return ret, nil
}

// RoundChangePrepared returns true if message is a RoundChange and prepared
func (msg *Message) RoundChangePrepared() bool {
	if msg.MsgType != RoundChangeMsgType {
		return false
	}

	return msg.DataRound != NoRound
}

// Encode returns a msg encoded bytes or error
func (msg *Message) Encode() ([]byte, error) {
	return msg.MarshalSSZ()
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	return msg.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (msg *Message) GetRoot() ([32]byte, error) {
	return msg.HashTreeRoot()
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (msg *Message) Validate() error {
	if len(msg.Identifier) == 0 {
		return errors.New("message identifier is invalid")
	}
	if _, err := msg.GetRoundChangeJustifications(); err != nil {
		return err
	}
	if _, err := msg.GetPrepareJustifications(); err != nil {
		return err
	}
	if msg.MsgType > RoundChangeMsgType {
		return errors.New("message type is invalid")
	}
	return nil
}

// Transforms a SignedSSVMessage into the same SignedSSVMessage with an empty Message.FullData
func SignedSSVMessageWithoutFullData(signedSSVMessage *types.SignedSSVMessage) (*types.SignedSSVMessage, error) {
	// Get Message
	message := &Message{}
	if err := message.Decode(signedSSVMessage.SSVMessage.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode Message")
	}
	// Make FullData empty
	message.FullData = []byte{}

	// Encode
	messageBytes, err := message.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode Message")
	}

	// New SSVMessage
	newSSVMessage := &types.SSVMessage{
		MsgType: signedSSVMessage.SSVMessage.MsgType,
		MsgID:   signedSSVMessage.SSVMessage.MsgID,
		Data:    messageBytes,
	}

	// New SignedSSVMessage
	return &types.SignedSSVMessage{
		OperatorID: signedSSVMessage.OperatorID,
		Signature:  signedSSVMessage.Signature,
		SSVMessage: newSSVMessage,
	}, nil
}
