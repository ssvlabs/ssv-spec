package qbft

import (
	"bytes"
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
func HasQuorum(share *types.Share, msgs []*SignedMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.GetSigners() {
			uniqueSigners[signer] = true
		}
	}
	return share.HasQuorum(len(uniqueSigners))
}

// HasPartialQuorum returns true if a unique set of signers has partial quorum
func HasPartialQuorum(share *types.Share, msgs []*SignedMessage) bool {
	uniqueSigners := make(map[types.OperatorID]bool)
	for _, msg := range msgs {
		for _, signer := range msg.GetSigners() {
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
	DataRound                Round
	RoundChangeJustification [][]byte `ssz-max:"13,65536"` // 2^16
	PrepareJustification     [][]byte `ssz-max:"13,65536"` // 2^16
}

func (msg *Message) GetRoundChangeJustifications() ([]*SignedMessage, error) {
	return unmarshalJustifications(msg.RoundChangeJustification)
}

func (msg *Message) GetPrepareJustifications() ([]*SignedMessage, error) {
	return unmarshalJustifications(msg.PrepareJustification)
}

func unmarshalJustifications(data [][]byte) ([]*SignedMessage, error) {
	ret := make([]*SignedMessage, len(data))
	for i, d := range data {
		sMsg := &SignedMessage{}
		if err := sMsg.UnmarshalSSZ(d); err != nil {
			return nil, err
		}
		ret[i] = sMsg
	}
	return ret, nil
}

func MarshalJustifications(msgs []*SignedMessage) ([][]byte, error) {
	ret := make([][]byte, len(msgs))
	for i, m := range msgs {
		d, err := m.WithoutFUllData().MarshalSSZ()
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
	if msg.MsgType > 5 {
		return errors.New("message type is invalid")
	}
	return nil
}

type SignedMessage struct {
	Signature types.Signature    `ssz-size:"96"`
	Signers   []types.OperatorID `ssz-max:"13"`
	Message   Message            // message for which this signature is for

	FullData []byte `ssz-max:"1073872896"` // 2^30+2^17
}

func (signedMsg *SignedMessage) GetSignature() types.Signature {
	return signedMsg.Signature
}
func (signedMsg *SignedMessage) GetSigners() []types.OperatorID {
	return signedMsg.Signers
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (signedMsg *SignedMessage) MatchedSigners(ids []types.OperatorID) bool {
	if len(signedMsg.Signers) != len(ids) {
		return false
	}

	for _, id := range signedMsg.Signers {
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
func (signedMsg *SignedMessage) CommonSigners(ids []types.OperatorID) bool {
	for _, id := range signedMsg.Signers {
		for _, id2 := range ids {
			if id == id2 {
				return true
			}
		}
	}
	return false
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (signedMsg *SignedMessage) Aggregate(sig types.MessageSignature) error {
	if signedMsg.CommonSigners(sig.GetSigners()) {
		return errors.New("duplicate signers")
	}

	r1, err := signedMsg.GetRoot()
	if err != nil {
		return errors.Wrap(err, "could not get signature root")
	}
	r2, err := sig.GetRoot()
	if err != nil {
		return errors.Wrap(err, "could not get signature root")
	}
	if !bytes.Equal(r1[:], r2[:]) {
		return errors.New("can't aggregate, roots not equal")
	}

	aggregated, err := signedMsg.Signature.Aggregate(sig.GetSignature())
	if err != nil {
		return errors.Wrap(err, "could not aggregate signatures")
	}
	signedMsg.Signature = aggregated
	signedMsg.Signers = append(signedMsg.Signers, sig.GetSigners()...)

	return nil
}

// Encode returns a msg encoded bytes or error
func (signedMsg *SignedMessage) Encode() ([]byte, error) {
	return signedMsg.MarshalSSZ()
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessage) Decode(data []byte) error {
	return signedMsg.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (signedMsg *SignedMessage) GetRoot() ([32]byte, error) {
	return signedMsg.Message.GetRoot()
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (signedMsg *SignedMessage) DeepCopy() *SignedMessage {
	ret := &SignedMessage{
		Signers:   make([]types.OperatorID, len(signedMsg.Signers)),
		Signature: make([]byte, len(signedMsg.Signature)),
	}
	copy(ret.Signers, signedMsg.Signers)
	copy(ret.Signature, signedMsg.Signature)

	ret.Message = Message{
		MsgType:    signedMsg.Message.MsgType,
		Height:     signedMsg.Message.Height,
		Round:      signedMsg.Message.Round,
		Identifier: make([]byte, len(signedMsg.Message.Identifier)),
		//Data:       make([]byte, len(signedMsg.Message.Data)),

		Root:                     signedMsg.Message.Root,
		DataRound:                signedMsg.Message.DataRound,
		PrepareJustification:     signedMsg.Message.PrepareJustification,
		RoundChangeJustification: signedMsg.Message.RoundChangeJustification,
	}
	copy(ret.Message.Identifier, signedMsg.Message.Identifier)
	//copy(ret.Message.Data, signedMsg.Message.Data)
	copy(ret.FullData, signedMsg.FullData)
	return ret
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (signedMsg *SignedMessage) Validate() error {
	if len(signedMsg.Signers) == 0 {
		return errors.New("message signers is empty")
	}

	// check unique signers
	signed := make(map[types.OperatorID]bool)
	for _, signer := range signedMsg.Signers {
		if signed[signer] {
			return errors.New("non unique signer")
		}
		if signer == 0 {
			return errors.New("signer ID 0 not allowed")
		}
		signed[signer] = true
	}

	return signedMsg.Message.Validate()
}

// WithoutFUllData returns SignedMessage without full data
func (signedMsg *SignedMessage) WithoutFUllData() *SignedMessage {
	return &SignedMessage{
		Signers:   signedMsg.Signers,
		Signature: signedMsg.Signature,
		Message:   signedMsg.Message,
	}
}
