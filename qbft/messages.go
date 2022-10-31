package qbft

import (
	"bytes"
	"encoding/json"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
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

type Data struct {
	Root   [32]byte `ssz-size:"32"`
	Source []byte   `ssz-max:"387173"`
}

// Message includes the full consensus input to be decided on, used for decided, proposal and round-change messages
type Message struct {
	Height Height
	Round  Round
	Input  *Data
	// PreparedRound an optional field used for round-change
	PreparedRound Round
}

// Encode returns a msg encoded bytes or error
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decode returns error if decoding failed
func (m *Message) Decode(data []byte) error {
	return json.Unmarshal(data, &m)
}

// Prepared returns error if decoding failed
func (m *Message) Prepared() bool {
	// TODO<olegshmuelov>: len(Message.Input) is always != 0 (this check already done in message validation),
	if m.PreparedRound != NoRound || len(m.Input.Source) != 0 {
		return true
	}
	return false
}

// GetRoot returns the root used for signing and verification
func (m *Message) GetRoot() ([]byte, error) {
	if m.Input == nil {
		return nil, errors.New("message input invalid")
	}
	return m.Input.Root[:], nil
	//marshaledRoot, err := m.Encode()
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not encode message")
	//}
	//ret := sha256.Sum256(marshaledRoot)
	//return ret[:], nil
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (m *Message) Validate(msgType types.MsgType) error {
	switch msgType {
	case types.ConsensusProposeMsgType, types.DecidedMsgType:
		if len(m.Input.Source) == 0 || m.Input.Root == [32]byte{} {
			return errors.New("message input data is invalid")
		}
	case types.ConsensusPrepareMsgType, types.ConsensusCommitMsgType:
		if m.Input.Root == [32]byte{} {
			return errors.New("message input data is invalid")
		}
	case types.ConsensusRoundChangeMsgType:
		// TODO<olegshmuelov>: the input data for round change can be nil?
	}

	return nil
}

// SignedMessage includes a signature over Message AND optional justification fields (not signed over)
type SignedMessage struct {
	Message   *Message
	Signers   []types.OperatorID `ssz-max:"13"`
	Signature types.Signature    `ssz-size:"96"`

	RoundChangeJustifications []*SignedMessage `ssz-max:"13"`
	ProposalJustifications    []*SignedMessage `ssz-max:"13"`
}

func (s *SignedMessage) GetSignature() types.Signature {
	return s.Signature
}

func (s *SignedMessage) GetSigners() []types.OperatorID {
	return s.Signers
}

// MatchedSigners returns true if the provided signer ids are equal to GetSignerIds() without order significance
func (s *SignedMessage) MatchedSigners(ids []types.OperatorID) bool {
	if len(s.Signers) != len(ids) {
		return false
	}

	for _, id := range s.Signers {
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
func (s *SignedMessage) CommonSigners(ids []types.OperatorID) bool {
	for _, id := range s.Signers {
		for _, id2 := range ids {
			if id == id2 {
				return true
			}
		}
	}
	return false
}

// Aggregate will aggregate the signed message if possible (unique signers, same digest, valid)
func (s *SignedMessage) Aggregate(sig types.MessageSignature) error {
	if s.CommonSigners(sig.GetSigners()) {
		return errors.New("duplicate signers")
	}

	r1, err := s.GetRoot()
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

	aggregated, err := s.Signature.Aggregate(sig.GetSignature())
	if err != nil {
		return errors.Wrap(err, "could not aggregate signatures")
	}
	s.Signature = aggregated
	s.Signers = append(s.Signers, sig.GetSigners()...)

	return nil
}

// Encode returns a msg encoded bytes or error
func (s *SignedMessage) Encode() ([]byte, error) {
	return json.Marshal(s)
}

// Decode returns error if decoding failed
func (s *SignedMessage) Decode(data []byte) error {
	return json.Unmarshal(data, &s)
}

// GetRoot returns the root used for signing and verification
func (s *SignedMessage) GetRoot() ([]byte, error) {
	return s.Message.GetRoot()
}

// DeepCopy returns a new instance of SignedMessage, deep copied
func (s *SignedMessage) DeepCopy(acceptedProposalData *Data) *SignedMessage {
	ret := &SignedMessage{
		Signers:   make([]types.OperatorID, len(s.Signers)),
		Signature: make([]byte, len(s.Signature)),
		// TODO<olegshmuelov>: handle justifications
		//RoundChangeJustifications: make([]*SignedMessage, len(sm.RoundChangeJustifications)),
		//ProposalJustifications:    make([]*SignedMessage, len(sm.ProposalJustifications)),
	}
	copy(ret.Signers, s.Signers)
	copy(ret.Signature, s.Signature)
	// TODO<olegshmuelov>: handle justifications
	//copy(ret.RoundChangeJustifications, sm.RoundChangeJustifications)
	//copy(ret.ProposalJustifications, sm.ProposalJustifications)

	ret.Message = &Message{
		Height: s.Message.Height,
		Round:  s.Message.Round,
		Input: &Data{
			Root:   acceptedProposalData.Root,
			Source: make([]byte, len(acceptedProposalData.Source)),
		},
		PreparedRound: s.Message.PreparedRound,
	}

	copy(ret.Message.Input.Source, acceptedProposalData.Source)
	return ret
}

// Validate returns error if msg validation doesn't pass.
// Msg validation checks the msg, it's variables for validity.
func (s *SignedMessage) Validate(msgType types.MsgType) error {
	if len(s.Signature) != 96 {
		return errors.New("message signature is invalid")
	}
	if len(s.Signers) == 0 {
		return errors.New("message signers is empty")
	}

	// check unique signers
	signed := make(map[types.OperatorID]bool)
	for _, signer := range s.Signers {
		if signed[signer] {
			return errors.New("non unique signer")
		}
		if signer == 0 {
			return errors.New("signer ID 0 not allowed")
		}
		signed[signer] = true
	}

	return s.Message.Validate(msgType)
}

func (s *SignedMessage) ToJustification() *SignedMessage {
	return &SignedMessage{
		Message: &Message{
			Height: s.Message.Height,
			Round:  s.Message.Round,
			Input: &Data{
				Root:   s.Message.Input.Root,
				Source: nil,
			},
			PreparedRound: s.Message.PreparedRound,
		},
		Signers:   s.Signers,
		Signature: s.Signature,
	}
}

type signedMessageSSZ struct {
	Message   *Message
	Signers   []uint64 `ssz-max:"13"`
	Signature []byte   `ssz-size:"96"`

	// TODO<olegshmuelov> calculate the real signedMessage size
	RoundChangeJustifications [][]byte `ssz-max:"13,400000"`
	ProposalJustifications    [][]byte `ssz-max:"13,400000"`
}

// SizeSSZ returns the ssz encoded size in bytes for the SignedMessage object
// TODO<olegshmuelov> the calculation of the signed msg size should be improved performance wise
func (s *SignedMessage) SizeSSZ() int {
	smSSZ, err := s.toSignedMessageSSZ()
	if err != nil {
		panic(err)
	}
	return smSSZ.SizeSSZ()
}

// MarshalSSZ ssz marshals the signedMessageSSZ object
func (s *SignedMessage) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the signedMessageSSZ object to a target array
func (s *SignedMessage) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	smSSZ, err := s.toSignedMessageSSZ()
	if err != nil {
		return nil, err
	}
	return smSSZ.MarshalSSZTo(buf)
}

func (s *SignedMessage) toSignedMessageSSZ() (*signedMessageSSZ, error) {
	signers := make([]uint64, len(s.Signers))
	for i, signer := range s.Signers {
		signers[i] = uint64(signer)
	}

	ss := &signedMessageSSZ{
		Message:                   s.Message,
		Signers:                   signers,
		Signature:                 s.Signature,
		RoundChangeJustifications: make([][]byte, len(s.RoundChangeJustifications)),
		ProposalJustifications:    make([][]byte, len(s.ProposalJustifications)),
	}
	for i, justification := range s.RoundChangeJustifications {
		marshalSSZ, err := justification.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		ss.RoundChangeJustifications[i] = marshalSSZ
	}
	for i, justification := range s.ProposalJustifications {
		marshalSSZ, err := justification.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		ss.ProposalJustifications[i] = marshalSSZ
	}

	return ss, nil
}

// UnmarshalSSZ ssz unmarshals the SignedMessage object
func (s *SignedMessage) UnmarshalSSZ(buf []byte) error {
	ss := &signedMessageSSZ{}
	err := ss.UnmarshalSSZ(buf)
	if err != nil {
		return err
	}
	s.Message = ss.Message
	s.Signers = make([]types.OperatorID, len(ss.Signers))
	for i, signer := range ss.Signers {
		s.Signers[i] = types.OperatorID(signer)
	}
	s.Signature = ss.Signature

	s.RoundChangeJustifications = make([]*SignedMessage, len(ss.RoundChangeJustifications))
	for i, justification := range ss.RoundChangeJustifications {
		j := &SignedMessage{}
		err := j.UnmarshalSSZ(justification)
		if err != nil {
			return err
		}
		s.RoundChangeJustifications[i] = j
	}

	s.ProposalJustifications = make([]*SignedMessage, len(ss.ProposalJustifications))
	for i, justification := range ss.ProposalJustifications {
		j := &SignedMessage{}
		err := j.UnmarshalSSZ(justification)
		if err != nil {
			return err
		}
		s.ProposalJustifications[i] = j
	}
	return nil
}
