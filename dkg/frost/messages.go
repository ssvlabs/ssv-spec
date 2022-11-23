package frost

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/dkg"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

type ProtocolMsg struct {
	Round              ProtocolRound       `json:"round,omitempty"`
	PreparationMessage *PreparationMessage `json:"preparation,omitempty"`
	Round1Message      *Round1Message      `json:"round1,omitempty"`
	Round2Message      *Round2Message      `json:"round2,omitempty"`
	BlameMessage       *BlameMessage       `json:"blame,omitempty"`
}

func (msg *ProtocolMsg) hasOnlyOneMsg() bool {
	var count = 0
	if msg.PreparationMessage != nil {
		count++
	}
	if msg.Round1Message != nil {
		count++
	}
	if msg.Round2Message != nil {
		count++
	}
	if msg.BlameMessage != nil {
		count++
	}
	return count == 1
}

func (msg *ProtocolMsg) msgMatchesRound() bool {
	switch msg.Round {
	case Preparation:
		return msg.PreparationMessage != nil
	case Round1:
		return msg.Round1Message != nil
	case Round2:
		return msg.Round2Message != nil
	case Blame:
		return msg.BlameMessage != nil
	default:
		return false
	}
}

func (msg *ProtocolMsg) Validate() error {
	if !msg.hasOnlyOneMsg() {
		return errors.New("need to contain one and only one message round")
	}
	if !msg.msgMatchesRound() {
		return errors.New("")
	}
	switch msg.Round {
	case Preparation:
		return msg.PreparationMessage.Validate()
	case Round1:
		return msg.Round1Message.Validate()
	case Round2:
		return msg.Round2Message.Validate()
	}
	return nil
}

// Encode returns a msg encoded bytes or error
func (msg *ProtocolMsg) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *ProtocolMsg) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

type PreparationMessage struct {
	SessionPk []byte
}

func (msg *PreparationMessage) Validate() error {
	_, err := ecies.NewPublicKeyFromBytes(msg.SessionPk)
	return err
}

// Encode returns a msg encoded bytes or error
func (msg *PreparationMessage) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *PreparationMessage) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

type Round1Message struct {
	// Commitment bytes representation of commitment points to pre-selected polynomials
	Commitment [][]byte
	// ProofS the S value of the Schnorr's proof
	ProofS []byte
	// ProofR the R value of the Schnorr's proof
	ProofR []byte
	// Shares the encrypted shares by operator
	Shares map[uint32][]byte
}

func (msg *Round1Message) Validate() error {
	var err error
	for _, bytes := range msg.Commitment {
		_, err = thisCurve.Point.FromAffineCompressed(bytes)
		if err != nil {
			return errors.Wrap(err, "invalid commitment")
		}
	}

	_, err = thisCurve.Scalar.SetBytes(msg.ProofS)
	if err != nil {
		return errors.Wrap(err, "invalid ProofS")
	}
	_, err = thisCurve.Scalar.SetBytes(msg.ProofR)
	if err != nil {
		return errors.Wrap(err, "invalid ProofR")
	}

	return nil
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

func (msg *Round2Message) Validate() error {
	var err error
	_, err = thisCurve.Point.FromAffineCompressed(msg.Vk)
	if err != nil {
		return errors.Wrap(err, "invalid vk")
	}
	_, err = thisCurve.Point.FromAffineCompressed(msg.VkShare)
	if err != nil {
		return errors.Wrap(err, "invalid vk share")
	}
	return nil
}

type BlameMessage struct {
	Type             BlameType
	TargetOperatorID uint32
	BlameData        [][]byte // SignedMessages received from the bad participant
	BlamerSessionSk  []byte
}

func (msg *BlameMessage) Validate() error {
	if len(msg.BlameData) < 1 {
		return errors.New("no blame data")
	}
	for _, datum := range msg.BlameData {
		signedMsg := &dkg.SignedMessage{}
		err := signedMsg.Decode(datum)
		if err != nil {
			return errors.Wrap(err, "contained data is not SignedMessage")
		}
	}
	return nil
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
	//// InvalidMessage refers to messages containing invalid values
	InvalidMessage
)

func (t BlameType) ToString() string {
	m := map[BlameType]string{
		InconsistentMessage: "Inconsistent Message",
		InvalidShare:        "Invalid Share",
		//FailedEcies:         "Failed Ecies",
		InvalidMessage: "Invalid Message",
	}
	return m[t]
}
