package keygen

import (
	"encoding/json"
	"errors"
)

type Round1Msg struct {
	Commitment []byte `json:"commitment"`
}

type Round2Msg struct {
	DeCommmitment [][]byte `json:"deCommitment"`
	BlindFactor   []byte   `json:"blindFactor"`
}

type Round3Msg struct {
	Share       []byte   `json:"share"`
}

type Round4Msg struct {
	Commitment        []byte `json:"commitment"`
	PubKey            []byte `json:"pubKey"`
	ChallengeResponse []byte `json:"challengeResponse"`
}

type MessageBody struct {
	Round1 *Round1Msg `json:"round1,omitempty"`
	Round2 *Round2Msg `json:"round2,omitempty"`
	Round3 *Round3Msg `json:"round3,omitempty"`
	Round4 *Round4Msg `json:"round4,omitempty"`
}

type Message struct {
	Sender   uint16      `json:"sender"`
	Receiver *uint16     `json:"receiver"`
	Body     MessageBody `json:"body"`
}

// Encode returns a msg encoded bytes or error
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decode returns error if decoding failed
func (m *Message) Decode(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Message) IsValid() bool {
	cnt := 0
	if m.Body.Round1 != nil {
		cnt += 1
	}
	if m.Body.Round2 != nil {
		cnt += 1
	}
	if m.Body.Round3 != nil {
		cnt += 1
	}
	if m.Body.Round4 != nil {
		cnt += 1
	}
	return cnt == 1
}

func (m *Message) GetRoundNumber() (int, error) {
	if m.Body.Round1 != nil {
		return 1, nil
	}
	if m.Body.Round2 != nil {
		return 2, nil
	}
	if m.Body.Round3 != nil {
		return 3, nil
	}
	if m.Body.Round4 != nil {
		return 4, nil
	}
	return 0, errors.New("invalid round")
}

type LocalKeyShare struct {
	Index           uint16   `json:"i"`
	Threshold       uint16   `json:"threshold"`
	ShareCount      uint16   `json:"share_count"`
	PublicKey       []byte   `json:"vk"`
	SecretShare     []byte   `json:"sk_i"`
	SharePublicKeys [][]byte `json:"vk_vec"`
}

type Messages = []*Message

// Encode returns a msg encoded bytes or error
func (msgs *Messages) Encode() ([]byte, error) {
	return json.Marshal(msgs)
}

// Decode returns error if decoding failed
func (msgs *Messages) Decode(data []byte) error {
	return json.Unmarshal(data, msgs)
}
