package qbft

import (
	"bytes"
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
	//marshaledRoot, err := msg.Encode()
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not encode message")
	//}
	//ret := sha256.Sum256(marshaledRoot)
	//return ret[:], nil
	cd := &types.ConsensusInput{}
	if err := cd.UnmarshalSSZ(msg.Input); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal consensus input ssz")
	}

	root, err := cd.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not hash tree consensus input root")
	}
	return root[:], nil
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
	return nil
}

func (msg *Message) ToMessageHeader() (*MessageHeader, error) {
	cd := &types.ConsensusInput{}
	if err := cd.UnmarshalSSZ(msg.Input); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal consensus input ssz")
	}

	root, err := cd.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not hash tree consensus input root")
	}

	return &MessageHeader{
		Height:        msg.Height,
		Round:         msg.Round,
		InputRoot:     root,
		PreparedRound: msg.PreparedRound,
	}, nil
}

func (msg *Message) GetHeaderInputRoot() ([32]byte, error) {
	cd := &types.ConsensusInput{}
	if err := cd.UnmarshalSSZ(msg.Input); err != nil {
		return [32]byte{}, errors.Wrap(err, "could not unmarshal consensus input ssz")
	}

	return cd.HashTreeRoot()
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
	return m.InputRoot[:], nil
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (m *MessageHeader) Validate() error {
	if m.InputRoot == [32]byte{} {
		return errors.New("message input data is invalid")
	}
	return nil
}

// SignedMessageHeader includes a signature over MessageHeader
type SignedMessageHeader struct {
	Message   *MessageHeader
	Signers   []types.OperatorID `ssz-max:"13"`
	Signature types.Signature    `ssz-size:"96"`
}

func (smh *SignedMessageHeader) GetSignature() types.Signature {
	return smh.Signature
}
func (smh *SignedMessageHeader) GetSigners() []types.OperatorID {
	return smh.Signers
}

// Encode returns a msg encoded bytes or error
func (smh *SignedMessageHeader) Encode() ([]byte, error) {
	return json.Marshal(smh)
}

// Decode returns error if decoding failed
func (smh *SignedMessageHeader) Decode(data []byte) error {
	return json.Unmarshal(data, &smh)
}

// GetRoot returns the root used for signing and verification
func (smh *SignedMessageHeader) GetRoot() ([]byte, error) {
	//TODO<olegshmuelov> implement
	//return smh.Message.GetRoot()
	return nil, nil
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (smh *SignedMessageHeader) MatchedSigners(ids []types.OperatorID) bool {
	if len(smh.Signers) != len(ids) {
		return false
	}

	for _, id := range smh.Signers {
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
func (smh *SignedMessageHeader) CommonSigners(ids []types.OperatorID) bool {
	for _, id := range smh.Signers {
		for _, id2 := range ids {
			if id == id2 {
				return true
			}
		}
	}
	return false
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (smh *SignedMessageHeader) Aggregate(sig types.MessageSignature) error {
	if smh.CommonSigners(sig.GetSigners()) {
		return errors.New("duplicate signers")
	}

	r1, err := smh.GetRoot()
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

	aggregated, err := smh.Signature.Aggregate(sig.GetSignature())
	if err != nil {
		return errors.Wrap(err, "could not aggregate signatures")
	}
	smh.Signature = aggregated
	smh.Signers = append(smh.Signers, sig.GetSigners()...)

	return nil
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (smh *SignedMessageHeader) DeepCopy(input []byte) *SignedMessage {
	ret := &SignedMessage{
		Signers:   make([]types.OperatorID, len(smh.Signers)),
		Signature: make([]byte, len(smh.Signature)),
		// TODO<olegshmuelov>: handle justifications
		//RoundChangeJustifications: make([]*SignedMessageHeader, len(sm.RoundChangeJustifications)),
		//ProposalJustifications:    make([]*SignedMessageHeader, len(sm.ProposalJustifications)),
	}
	copy(ret.Signers, smh.Signers)
	copy(ret.Signature, smh.Signature)
	// TODO<olegshmuelov>: handle justifications
	//copy(ret.RoundChangeJustifications, sm.RoundChangeJustifications)
	//copy(ret.ProposalJustifications, sm.ProposalJustifications)

	ret.Message = &Message{
		Height:        smh.Message.Height,
		Round:         smh.Message.Round,
		Input:         make([]byte, len(input)),
		PreparedRound: smh.Message.PreparedRound,
	}

	copy(ret.Message.Input, input)
	return ret
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (smh *SignedMessageHeader) Validate() error {
	if len(smh.Signature) != 96 {
		return errors.New("message signature is invalid")
	}
	if len(smh.Signers) == 0 {
		return errors.New("message signers is empty")
	}

	// check unique signers
	signed := make(map[types.OperatorID]bool)
	for _, signer := range smh.Signers {
		if signed[signer] {
			return errors.New("non unique signer")
		}
		signed[signer] = true
	}

	return smh.Message.Validate()
}

func (sm *SignedMessage) GetSignature() types.Signature {
	return sm.Signature
}
func (sm *SignedMessage) GetSigners() []types.OperatorID {
	return sm.Signers
}

func (sm *SignedMessage) ToSignedMessageHeader() (*SignedMessageHeader, error) {
	header, err := sm.Message.ToMessageHeader()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert to header")
	}

	return &SignedMessageHeader{
		Message:   header,
		Signers:   sm.Signers,
		Signature: sm.Signature,
	}, nil
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
		return errors.New("duplicate signers")
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
		Signers:                   make([]types.OperatorID, len(sm.Signers)),
		Signature:                 make([]byte, len(sm.Signature)),
		RoundChangeJustifications: make([]*SignedMessageHeader, len(sm.RoundChangeJustifications)),
		ProposalJustifications:    make([]*SignedMessageHeader, len(sm.ProposalJustifications)),
	}
	copy(ret.Signers, sm.Signers)
	copy(ret.Signature, sm.Signature)
	copy(ret.RoundChangeJustifications, sm.RoundChangeJustifications)
	copy(ret.ProposalJustifications, sm.ProposalJustifications)

	ret.Message = &Message{
		Height:        sm.Message.Height,
		Round:         sm.Message.Round,
		Input:         make([]byte, len(sm.Message.Input)),
		PreparedRound: sm.Message.PreparedRound,
	}

	copy(ret.Message.Input, sm.Message.Input)
	return ret
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (sm *SignedMessage) Validate(msgType types.MsgType) error {
	if len(sm.Signature) != 96 {
		return errors.New("message signature is invalid")
	}
	if len(sm.Signers) == 0 {
		return errors.New("message signers is empty")
	}

	// check unique signers
	signed := make(map[types.OperatorID]bool)
	for _, signer := range sm.Signers {
		if signed[signer] {
			return errors.New("non unique signer")
		}
		signed[signer] = true
	}

	return sm.Message.Validate(msgType)
}
