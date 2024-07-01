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
func HasQuorum(share *types.CommitteeMember, msgs []*ProcessingMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.SignedMessage.OperatorIDs {
			uniqueSigners[signer] = true
		}
	}
	return share.HasQuorum(len(uniqueSigners))
}

// HasPartialQuorum returns true if a unique set of signers has partial quorum
func HasPartialQuorum(share *types.CommitteeMember, msgs []*ProcessingMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.SignedMessage.OperatorIDs {
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
	RoundChangeJustification [][]byte `ssz-max:"13,51852"`
	PrepareJustification     [][]byte `ssz-max:"13,3700"`
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

func Sign(msg *Message, operatorID types.OperatorID, operatorSigner types.OperatorSigner) (*types.SignedSSVMessage, error) {

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

	sig, err := operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign SSVMessage")
	}

	signedSSVMessage := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{operatorID},
		SSVMessage:  ssvMsg,
	}

	return signedSSVMessage, nil
}

// ProcessingMessage stores the network-exchanged signed message and the decoded SignedMessage.SSVMessage.Data (qbft message)
// The signed message is used at the qbft level since it's used for qbft justifications
type ProcessingMessage struct {
	SignedMessage *types.SignedSSVMessage
	QBFTMessage   *Message
}

// NewProcessingMessage creates a ProcessingMessage with the decoded qbft message
func NewProcessingMessage(signedMessage *types.SignedSSVMessage) (*ProcessingMessage, error) {
	msg := &Message{}
	err := msg.Decode(signedMessage.SSVMessage.Data)
	if err != nil {
		return nil, err
	}
	return &ProcessingMessage{
		SignedMessage: signedMessage,
		QBFTMessage:   msg,
	}, nil
}

// Validate checks the signed message and qbft message validation
func (msg *ProcessingMessage) Validate() error {
	if err := msg.SignedMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}
	if err := msg.QBFTMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid Message")
	}
	return nil
}
