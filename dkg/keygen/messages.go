package keygen

import "errors"

type Round1Msg struct {
	Commitment []byte `json:"commitment"`
}

type Round2Msg struct {
	YI          []byte `json:"yI"`
	BlindFactor []byte `json:"blindFactor"`
}

type Round3Msg struct {
	Commitments [][]byte `json:"commitments"`
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

type Messages = []*Message

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
