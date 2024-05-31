package qbft

import (
	"crypto/sha256"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/types"
)

// HashDataRoot hashes input data to root
func HashDataRoot(data []byte) ([32]byte, error) {
	ret := sha256.Sum256(data)
	return ret, nil
}

// HasQuorum returns true if a unique set of signers has quorum
func HasQuorum(share *types.SharedValidator, msgs []*types.SignedSSVMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.GetOperatorIDs() {
			uniqueSigners[signer] = true
		}
	}
	return share.HasQuorum(len(uniqueSigners))
}

// HasPartialQuorum returns true if a unique set of signers has partial quorum
func HasPartialQuorum(share *types.SharedValidator, msgs []*types.SignedSSVMessage) bool {
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
}

// Creates a Message object from bytes
func DecodeMessage(data []byte) (*Message, error) {
	ret := &Message{}
	err := ret.Decode(data)
	return ret, err
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
		d, err := m.WithoutFullData().MarshalSSZ()
		if err != nil {
			return nil, err
		}
		ret[i] = d
	}
	return ret, nil
}

func MessageToSignedSSVMessageWithFullData(msg *Message, operatorID types.OperatorID, operatorSigner types.OperatorSigner, fullData []byte) (*types.SignedSSVMessage, error) {
	signedSSVMessage, err := MessageToSignedSSVMessage(msg, operatorID, operatorSigner)
	if err != nil {
		return nil, err
	}
	signedSSVMessage.FullData = fullData
	return signedSSVMessage, nil
}

func MessageToSignedSSVMessage(msg *Message, operatorID types.OperatorID, operatorSigner types.OperatorSigner) (*types.SignedSSVMessage, error) {

	byts, err := msg.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode message")
	}

	msgID := types.MessageID{}
	copy(msgID[:], msg.Identifier)

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    byts,
	}

	signedSSVMessage, err := types.SSVMessageToSignedSSVMessage(ssvMsg, operatorID, operatorSigner.SignSSVMessage)
	if err != nil {
		return nil, errors.Wrap(err, "could not create SignedSSVMessage from SSVMessage")
	}
	return signedSSVMessage, nil
}
