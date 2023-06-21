package keysign

import (
	"encoding/json"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/common"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

type ProtocolMsg struct {
	Round              common.ProtocolRound `json:"round,omitempty"`
	PreparationMessage *PreparationMessage  `json:"preparation,omitempty"`
	Round1Message      *Round1Message       `json:"round1,omitempty"`
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
	return count == 1
}

func (msg *ProtocolMsg) msgMatchesRound() bool {
	switch msg.Round {
	case common.Preparation:
		return msg.PreparationMessage != nil
	case common.Round1:
		return msg.Round1Message != nil
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

// Encode returns a msg encoded bytes or error
func (msg *ProtocolMsg) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

// Decode returns error if decoding failed
func (msg *ProtocolMsg) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}

type PreparationMessage struct {
	PartialSignature []byte
}

func (msg *PreparationMessage) Validate() error {
	return nil
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
	ReconstructedSignature []byte
}

func (msg *Round1Message) Validate() error {
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

type TimeoutMessage struct {
	Round common.ProtocolRound
}
