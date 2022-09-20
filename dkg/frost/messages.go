package frost

import "encoding/json"

type ProtocolMsg struct {
	Round              DKGRound            `json:"round,omitempty"`
	PreparationMessage *PreparationMessage `json:"preparation,omitempty"`
	Round1Message      *Round1Message      `json:"round1,omitempty"`
	Round2Message      *Round2Message      `json:"round2,omitempty"`
	BlameMessage       *BlameMessage       `json:"blame,omitempty"`
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

type Round2Message struct {
	Vk      []byte
	VkShare []byte
}

type BlameMessage struct {
	TargetOperatorID uint64
	BlameData        [][]byte // SignedMessages received from the bad participant
	BlamerSessionSk  []byte
}

// Encode returns a msg encoded bytes or error
func (msg *ProtocolMsg) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *ProtocolMsg) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
