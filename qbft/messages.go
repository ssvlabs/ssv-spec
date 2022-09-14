package qbft

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

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

// HasQuorumHeaders returns true if a unique set of signers has quorum
func HasQuorumHeaders(share *types.Share, msgs []*SignedMessageHeader) bool {
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

/*func (d *RoundChangeData) Prepared() bool {
	if d.PreparedRound != NoRound || len(d.PreparedValue) != 0 {
		return true
	}
	return false
}*/

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
/*func (d *RoundChangeData) Validate() error {
	if d.Prepared() {
		if len(d.PreparedValue) == 0 {
			return errors.New("round change prepared value invalid")
		}
		if len(d.RoundChangeJustification) == 0 {
			return errors.New("round change justification invalid")
		}
		// TODO - should next proposal data be equal to prepared value?
	}
	return nil
}*/

// Message includes the full consensus input to be decided on, used for decided, proposal and round-change messages
type Message struct {
	Height Height
	Round  Round
	Input  []byte `ssz-max:"387173"`
	// PreparedRound an optional field used for round-change
	PreparedRound Round
}

// Encode returns a msg encoded bytes or error
func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Message) Decode(data []byte) error {
	return json.Unmarshal(data, &msg)
}

// GetRoot returns the root used for signing and verification
func (msg *Message) GetRoot() ([]byte, error) {
	marshaledRoot, err := msg.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode message")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (msg *Message) Validate(msgType types.MsgType) error {
	//if len(msg.Identifier) == 0 {
	//	return errors.New("message identifier is invalid")
	//}
	if len(msg.Input) == 0 && msgType != types.ConsensusRoundChangeMsgType {
		return errors.New("message input data is invalid")
	}
	//if msg.MsgType > 5 {
	//	return errors.New("message type is invalid")
	//}
	return nil
}

func (msg *Message) ToMessageHeader() (*MessageHeader, error) {
	// TODO<olegshmuelov>: implement HashTreeRoot ssz
	//r, err := msg.Input.HashTreeRoot()
	//if err != nil {
	//	return &MessageHeader{}, errors.Wrap(err, "failed to get input root")
	//}
	return &MessageHeader{
		Height: msg.Height,
		Round:  msg.Round,
		// TODO<olegshmuelov>: implement HashTreeRoot ssz
		InputRoot:     [32]byte{},
		PreparedRound: msg.PreparedRound,
	}, nil
}

// SignedMessage includes a signature over Message AND optional justification fields (not signed over)
type SignedMessage struct {
	Message   *Message
	Signers   []types.OperatorID `ssz-max:"13"`
	Signature types.Signature    `ssz-size:"96"`

	RoundChangeJustifications []*SignedMessageHeader `ssz-max:"13"`
	ProposalJustifications    []*SignedMessageHeader `ssz-max:"13"`
}

// MessageHeader includes just the root of the input to be decided on (to save space), used for prepare, commit and justification messages
type MessageHeader struct {
	Height        Height
	Round         Round
	InputRoot     [32]byte `ssz-size:"32"`
	PreparedRound Round
}

// Encode returns a msg encoded bytes or error
func (m *MessageHeader) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decode returns error if decoding failed
func (m *MessageHeader) Decode(data []byte) error {
	return json.Unmarshal(data, &m)
}

func (m *MessageHeader) GetRoot() ([]byte, error) {
	marshaledRoot, err := m.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode message")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// SignedMessageHeader includes a signature over MessageHeader
type SignedMessageHeader struct {
	Message   *MessageHeader
	Signers   []types.OperatorID `ssz-max:"13"`
	Signature types.Signature    `ssz-size:"96"`
}

func (signedMsg *SignedMessageHeader) GetSignature() types.Signature {
	return signedMsg.Signature
}
func (signedMsg *SignedMessageHeader) GetSigners() []types.OperatorID {
	return signedMsg.Signers
}

// Encode returns a msg encoded bytes or error
func (signedMsg *SignedMessageHeader) Encode() ([]byte, error) {
	return json.Marshal(signedMsg)
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessageHeader) Decode(data []byte) error {
	return json.Unmarshal(data, &signedMsg)
}

// GetRoot returns the root used for signing and verification
func (signedMsg *SignedMessageHeader) GetRoot() ([]byte, error) {
	//TODO<olegshmuelov> implement
	//return signedMsg.Message.GetRoot()
	return nil, nil
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (signedMsg *SignedMessageHeader) MatchedSigners(ids []types.OperatorID) bool {
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

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (signedMsg *SignedMessageHeader) Aggregate(sig types.MessageSignature) error {
	// TODO<olegshmuelov> implement
	//if signedMsg.CommonSigners(sig.GetSigners()) {
	//	return errors.New("duplicate signers")
	//}

	r1, err := signedMsg.GetRoot()
	if err != nil {
		return errors.Wrap(err, "could not get signature root")
	}
	r2, err := sig.GetRoot()
	if err != nil {
		return errors.Wrap(err, "could not get signature root")
	}
	if !bytes.Equal(r1, r2) {
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

func (signedMsg *SignedMessage) GetSignature() types.Signature {
	return signedMsg.Signature
}
func (signedMsg *SignedMessage) GetSigners() []types.OperatorID {
	return signedMsg.Signers
}

func (signedMsg *SignedMessage) ToSignedMessageHeader() (*SignedMessageHeader, error) {
	header, err := signedMsg.Message.ToMessageHeader()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert to header")
	}

	return &SignedMessageHeader{
		Message:   header,
		Signers:   signedMsg.Signers,
		Signature: signedMsg.Signature,
	}, nil
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
	if !bytes.Equal(r1, r2) {
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
	return json.Marshal(signedMsg)
}

// Decode returns error if decoding failed
func (signedMsg *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &signedMsg)
}

// GetRoot returns the root used for signing and verification
func (signedMsg *SignedMessage) GetRoot() ([]byte, error) {
	return signedMsg.Message.GetRoot()
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (signedMsg *SignedMessage) DeepCopy() *SignedMessage {
	ret := &SignedMessage{
		Signers:                   make([]types.OperatorID, len(signedMsg.Signers)),
		Signature:                 make([]byte, len(signedMsg.Signature)),
		RoundChangeJustifications: make([]*SignedMessageHeader, len(signedMsg.RoundChangeJustifications)),
		ProposalJustifications:    make([]*SignedMessageHeader, len(signedMsg.ProposalJustifications)),
	}
	copy(ret.Signers, signedMsg.Signers)
	copy(ret.Signature, signedMsg.Signature)
	copy(ret.RoundChangeJustifications, signedMsg.RoundChangeJustifications)
	copy(ret.ProposalJustifications, signedMsg.ProposalJustifications)

	ret.Message = &Message{
		Height:        signedMsg.Message.Height,
		Round:         signedMsg.Message.Round,
		Input:         make([]byte, len(signedMsg.Message.Input)),
		PreparedRound: signedMsg.Message.PreparedRound,
	}

	copy(ret.Message.Input, signedMsg.Message.Input)
	return ret
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (signedMsg *SignedMessage) Validate(msgType types.MsgType) error {
	if len(signedMsg.Signature) != 96 {
		return errors.New("message signature is invalid")
	}
	if len(signedMsg.Signers) == 0 {
		return errors.New("message signers is empty")
	}

	// check unique signers
	signed := make(map[types.OperatorID]bool)
	for _, signer := range signedMsg.Signers {
		if signed[signer] {
			return errors.New("non unique signer")
		}
		signed[signer] = true
	}

	return signedMsg.Message.Validate(msgType)
}
