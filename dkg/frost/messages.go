package frost

import "encoding/json"

type ProtocolMsg struct {
	Round              DKGRound            `json:"round,omitempty"`
	PreparationMessage *PreparationMessage `json:"preparation,omitempty"`
	Round1Message      *Round1Message      `json:"round1,omitempty"`
	Round2Message      *Round2Message      `json:"round2,omitempty"`
	BlameMessage       *BlameMessage       `json:"blame,omitempty"`
}

func (msg *ProtocolMsg) validate() bool {
	var messageExists bool
	switch msg.Round {
	case Preparation:
		messageExists = msg.PreparationMessage != nil
	case Round1:
		messageExists = msg.Round1Message != nil
	case Round2:
		messageExists = msg.Round2Message != nil
	case Blame:
		messageExists = msg.BlameMessage != nil
	default:
		return false
	}
	return messageExists
}

type PreparationMessage struct {
	SessionPk []byte
}

type Round1Message struct {
	Commitment [][]byte
	ProofS     []byte
	ProofR     []byte
	Shares     map[uint32][]byte
}

// Encode returns a msg encoded bytes or error
func (msg *Round1Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *Round1Message) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

type Round2Message struct {
	Vk      []byte
	VkShare []byte
}

type BlameMessage struct {
	Type             BlameType
	TargetOperatorID uint32
	BlameData        [][]byte // SignedMessages received from the bad participant
	BlamerSessionSk  []byte
}

// Encode returns a msg encoded bytes or error
func (msg *BlameMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *BlameMessage) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

type BlameType int

const (
	// InconsistentMessage refers to an operator sending multiple messages for same round
	InconsistentMessage BlameType = iota
	// InvalidShare refers to an operator sending invalid share
	InvalidShare
)

// Encode returns a msg encoded bytes or error
func (msg *ProtocolMsg) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *ProtocolMsg) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
