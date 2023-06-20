package frost

import (
	"encoding/json"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"
	ecies "github.com/ecies/go/v2"
	"github.com/pkg/errors"
)

type ProtocolMsg struct {
	Round              common.ProtocolRound `json:"round,omitempty"`
	PreparationMessage *PreparationMessage  `json:"preparation,omitempty"`
	Round1Message      *Round1Message       `json:"round1,omitempty"`
	Round2Message      *Round2Message       `json:"round2,omitempty"`
	BlameMessage       *BlameMessage        `json:"blame,omitempty"`
	TimeoutMessage     *TimeoutMessage      `json:"timeout,omitempty"`
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
	case common.Preparation:
		return msg.PreparationMessage != nil
	case common.Round1:
		return msg.Round1Message != nil
	case common.Round2:
		return msg.Round2Message != nil
	case common.Blame:
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
	case common.Preparation:
		return msg.PreparationMessage.Validate()
	case common.Round1:
		return msg.Round1Message.Validate()
	case common.Round2:
		return msg.Round2Message.Validate()
	}
	return nil
}

func (msg *ProtocolMsg) ToSignedMessage(id dkg.RequestID, operatorID types.OperatorID, storage dkg.Storage, signer types.DKGSigner) (*dkg.SignedMessage, error) {
	msgBytes, err := msg.Encode()
	if err != nil {
		return nil, err
	}

	bcastMessage := &dkg.SignedMessage{
		Message: &dkg.Message{
			MsgType:    dkg.ProtocolMsgType,
			Identifier: id,
			Data:       msgBytes,
		},
		Signer: operatorID,
	}

	exist, operator, err := storage.GetDKGOperator(operatorID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.Errorf("operator with id %d not found", operatorID)
	}

	sig, err := signer.SignDKGOutput(bcastMessage, operator.ETHAddress)
	if err != nil {
		return nil, err
	}
	bcastMessage.Signature = sig
	return bcastMessage, nil
}

func (msg *ProtocolMsg) GetMessage(msgContainer common.IMsgContainer, round common.ProtocolRound, operatorID uint32) (interface{}, error) {
	signedMsg, err := msgContainer.GetSignedMsg(round, operatorID)
	if err != nil {
		return nil, err
	}
	if err := msg.Decode(signedMsg.Message.Data); err != nil {
		return nil, err
	}
	switch round {
	case common.Preparation:
		return msg.PreparationMessage, nil
	case common.Round1:
		return msg.Round1Message, nil
	case common.Round2:
		return msg.Round2Message, nil
	case common.Blame:
		return msg.BlameMessage, nil
	default:
		return nil, dkg.ErrInvalidRound{}
	}
}

func GetPreparationMsg(msgContainer common.IMsgContainer, operatorID uint32) (*PreparationMessage, error) {
	protocolMsg := new(ProtocolMsg)
	msg, err := protocolMsg.GetMessage(msgContainer, common.Preparation, operatorID)
	if err != nil {
		return nil, err
	}
	parsedMsg, _ := msg.(*PreparationMessage)
	if parsedMsg == nil {
		return nil, common.ErrMsgNil{Round: common.Preparation, OperatorID: operatorID}
	}
	return parsedMsg, nil
}

func GetRound1Msg(msgContainer common.IMsgContainer, operatorID uint32) (*Round1Message, error) {
	protocolMsg := new(ProtocolMsg)
	msg, err := protocolMsg.GetMessage(msgContainer, common.Round1, operatorID)
	if err != nil {
		return nil, err
	}
	parsedMsg, _ := msg.(*Round1Message)
	if parsedMsg == nil {
		return nil, common.ErrMsgNil{Round: common.Round1, OperatorID: operatorID}
	}
	return parsedMsg, nil
}

func GetRound2Msg(msgContainer common.IMsgContainer, operatorID uint32) (*Round2Message, error) {
	protocolMsg := new(ProtocolMsg)
	msg, err := protocolMsg.GetMessage(msgContainer, common.Round2, operatorID)
	if err != nil {
		return nil, err
	}
	parsedMsg, _ := msg.(*Round2Message)
	if parsedMsg == nil {
		return nil, common.ErrMsgNil{Round: common.Round2, OperatorID: operatorID}
	}
	return parsedMsg, nil
}

func GetBlameMsg(msgContainer common.IMsgContainer, operatorID uint32) (*BlameMessage, error) {
	protocolMsg := new(ProtocolMsg)
	msg, err := protocolMsg.GetMessage(msgContainer, common.Blame, operatorID)
	if err != nil {
		return nil, err
	}
	parsedMsg, _ := msg.(*BlameMessage)
	if parsedMsg == nil {
		return nil, common.ErrMsgNil{Round: common.Blame, OperatorID: operatorID}
	}
	return parsedMsg, nil
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
	// InvalidMessage refers to messages containing invalid values
	InvalidMessage
)

func (t BlameType) ToString() string {
	m := map[BlameType]string{
		InconsistentMessage: "Inconsistent Message",
		InvalidShare:        "Invalid Share",
		InvalidMessage:      "Invalid Message",
	}
	return m[t]
}

type TimeoutMessage struct {
	Round common.ProtocolRound
}
