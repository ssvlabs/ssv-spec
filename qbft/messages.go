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

type MessageType int

const (
	ProposalMsgType MessageType = iota
	PrepareMsgType
	CommitMsgType
	RoundChangeMsgType
)

type ProposalData struct {
	Data                     []byte           `ssz-max:"1161"` // TODO(olegshmuelov): check if the ssz-max is correct
	RoundChangeJustification []*SignedMessage `ssz-max:"10"`   // TODO(olegshmuelov): check if the ssz-max is correct
	PrepareJustification     []*SignedMessage `ssz-max:"10"`   // TODO(olegshmuelov): check if the ssz-max is correct
}

// Encode returns a msg encoded bytes or error
func (d *ProposalData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *ProposalData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *ProposalData) Validate() error {
	if len(d.Data) == 0 {
		return errors.New("ProposalData data is invalid")
	}
	return nil
}

type PrepareData struct {
	Data []byte
}

// Encode returns a msg encoded bytes or error
func (d *PrepareData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *PrepareData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *PrepareData) Validate() error {
	if len(d.Data) == 0 {
		return errors.New("PrepareData data is invalid")
	}
	return nil
}

type CommitData struct {
	Data []byte
}

// Encode returns a msg encoded bytes or error
func (d *CommitData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *CommitData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *CommitData) Validate() error {
	if len(d.Data) == 0 {
		return errors.New("CommitData data is invalid")
	}
	return nil
}

type RoundChangeData struct {
	PreparedValue            []byte
	PreparedRound            Round
	NextProposalData         []byte
	RoundChangeJustification []*SignedMessage
}

func (d *RoundChangeData) Prepared() bool {
	if d.PreparedRound != NoRound || len(d.PreparedValue) != 0 {
		return true
	}
	return false
}

// Encode returns a msg encoded bytes or error
func (d *RoundChangeData) Encode() ([]byte, error) {
	return json.Marshal(d)
}

// Decode returns error if decoding failed
func (d *RoundChangeData) Decode(data []byte) error {
	return json.Unmarshal(data, &d)
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (d *RoundChangeData) Validate() error {
	if d.Prepared() {
		if len(d.PreparedValue) == 0 {
			return errors.New("round change prepared value invalid")
		}
		if len(d.RoundChangeJustification) == 0 {
			return errors.New("round change justification invalid")
		}
		// TODO - should next proposal data be equal to prepared value?
	}

	if len(d.NextProposalData) == 0 {
		return errors.New("round change next value invalid")
	}
	return nil
}

type Message struct {
	MsgType    MessageType
	Height     Height // QBFT instance Height
	Round      Round  // QBFT round for which the msg is for
	Identifier []byte `ssz-max:"96"`   // instance Identifier this msg belongs to // TODO(olegshmuelov): check if the ssz-max is correct
	Data       []byte `ssz-max:"2000"` // TODO(olegshmuelov): check if the ssz-max is correct
}

type messageSSZ struct {
	MsgType    uint8
	Height     uint64
	Round      Round
	Identifier []byte `ssz-max:"96"`   // TODO(olegshmuelov): check if the ssz-max is correct
	Data       []byte `ssz-max:"2000"` // TODO(olegshmuelov): check if the ssz-max is correct
}

// GetProposalData returns proposal specific data
func (m *Message) GetProposalData() (*ProposalData, error) {
	ret := &ProposalData{}
	if err := ret.Decode(m.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode proposal data from message")
	}
	return ret, nil
}

// GetPrepareData returns prepare specific data
func (m *Message) GetPrepareData() (*PrepareData, error) {
	ret := &PrepareData{}
	if err := ret.Decode(m.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode prepare data from message")
	}
	return ret, nil
}

// GetCommitData returns commit specific data
func (m *Message) GetCommitData() (*CommitData, error) {
	ret := &CommitData{}
	if err := ret.Decode(m.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode commit data from message")
	}
	return ret, nil
}

// GetRoundChangeData returns round change specific data
func (m *Message) GetRoundChangeData() (*RoundChangeData, error) {
	ret := &RoundChangeData{}
	if err := ret.Decode(m.Data); err != nil {
		return nil, errors.Wrap(err, "could not decode round change data from message")
	}
	return ret, nil
}

// Encode returns a msg encoded bytes or error
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decode returns error if decoding failed
func (m *Message) Decode(data []byte) error {
	return json.Unmarshal(data, &m)
}

// GetRoot returns the root used for signing and verification
func (m *Message) GetRoot() ([]byte, error) {
	marshaledRoot, err := m.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode message")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (m *Message) Validate() error {
	if len(m.Identifier) == 0 {
		return errors.New("message identifier is invalid")
	}
	if len(m.Data) == 0 {
		return errors.New("message data is invalid")
	}
	if m.MsgType > 5 {
		return errors.New("message type is invalid")
	}
	return nil
}

type SignedMessage struct {
	Signature types.Signature    `ssz-max:"100"`
	Signers   []types.OperatorID `ssz-max:"7"`
	Message   *Message           // message for which this signature is for
}

type signedMessageSSZ struct {
	Signature types.Signature `ssz-max:"100"`
	Signers   []uint64        `ssz-max:"7"`
	Message   *Message        // message for which this signature is for
}

func (sm *SignedMessage) GetSignature() types.Signature {
	return sm.Signature
}
func (sm *SignedMessage) GetSigners() []types.OperatorID {
	return sm.Signers
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (sm *SignedMessage) MatchedSigners(ids []types.OperatorID) bool {
	if len(sm.Signers) != len(ids) {
		return false
	}

	for _, id := range sm.Signers {
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
func (sm *SignedMessage) CommonSigners(ids []types.OperatorID) bool {
	for _, id := range sm.Signers {
		for _, id2 := range ids {
			if id == id2 {
				return true
			}
		}
	}
	return false
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (sm *SignedMessage) Aggregate(sig types.MessageSignature) error {
	if sm.CommonSigners(sig.GetSigners()) {
		return errors.New("can't aggregate 2 signed messages with mutual signers")
	}

	r1, err := sm.GetRoot()
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

	aggregated, err := sm.Signature.Aggregate(sig.GetSignature())
	if err != nil {
		return errors.Wrap(err, "could not aggregate signatures")
	}
	sm.Signature = aggregated
	sm.Signers = append(sm.Signers, sig.GetSigners()...)

	return nil
}

// Encode returns a msg encoded bytes or error
func (sm *SignedMessage) Encode() ([]byte, error) {
	return json.Marshal(sm)
}

// Decode returns error if decoding failed
func (sm *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &sm)
}

// GetRoot returns the root used for signing and verification
func (sm *SignedMessage) GetRoot() ([]byte, error) {
	return sm.Message.GetRoot()
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (sm *SignedMessage) DeepCopy() *SignedMessage {
	ret := &SignedMessage{
		Signers:   make([]types.OperatorID, len(sm.Signers)),
		Signature: make([]byte, len(sm.Signature)),
	}
	copy(ret.Signers, sm.Signers)
	copy(ret.Signature, sm.Signature)

	ret.Message = &Message{
		MsgType:    sm.Message.MsgType,
		Height:     sm.Message.Height,
		Round:      sm.Message.Round,
		Identifier: make([]byte, len(sm.Message.Identifier)),
		Data:       make([]byte, len(sm.Message.Data)),
	}
	copy(ret.Message.Identifier, sm.Message.Identifier)
	copy(ret.Message.Data, sm.Message.Data)
	return ret
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (sm *SignedMessage) Validate() error {
	if len(sm.Signature) != 96 {
		return errors.New("message signature is invalid")
	}
	if len(sm.Signers) == 0 {
		return errors.New("message signers is empty")
	}
	return sm.Message.Validate()
}

type DecidedMessage struct {
	SignedMessage *SignedMessage
}

// Encode returns a msg encoded bytes or error
func (msg *DecidedMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *DecidedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &msg)
}
