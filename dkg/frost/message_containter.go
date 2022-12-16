package frost

import (
	"fmt"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/pkg/errors"
)

type MsgContainer struct {
	msgs map[ProtocolRound]map[uint32]*dkg.SignedMessage
}

func newMsgContainer() *MsgContainer {
	m := make(map[ProtocolRound]map[uint32]*dkg.SignedMessage)
	for _, round := range rounds {
		m[round] = make(map[uint32]*dkg.SignedMessage)
	}
	return &MsgContainer{msgs: m}
}

func (msgContainer *MsgContainer) SaveMsg(round ProtocolRound, msg *dkg.SignedMessage) (existingMessage *dkg.SignedMessage, err error) {
	existingMessage, exists := msgContainer.msgs[round][uint32(msg.Signer)]
	if exists {
		return existingMessage, errors.New("msg already exists")
	}
	msgContainer.msgs[round][uint32(msg.Signer)] = msg
	return nil, nil
}

func (msgContainer *MsgContainer) GetSignedMessage(round ProtocolRound, operatorID uint32) (*dkg.SignedMessage, error) {
	signedMsg, exist := msgContainer.msgs[round][operatorID]
	if !exist {
		return nil, ErrMsgNotFound{round: round, operatorID: operatorID}
	}
	return signedMsg, nil
}

func (msgContainer *MsgContainer) GetPreparationMsg(operatorID uint32) (*PreparationMessage, error) {
	msg, err := msgContainer.GetMessage(Preparation, operatorID)
	if err != nil {
		return nil, err
	}
	prepMsg, _ := msg.(*PreparationMessage)
	if prepMsg == nil {
		return nil, ErrMsgNil{round: Preparation, operatorID: operatorID}
	}
	return prepMsg, nil
}

func (msgContainer *MsgContainer) GetRound1Msg(operatorID uint32) (*Round1Message, error) {
	msg, err := msgContainer.GetMessage(Round1, operatorID)
	if err != nil {
		return nil, err
	}
	prepMsg, _ := msg.(*Round1Message)
	if prepMsg == nil {
		return nil, ErrMsgNil{round: Round1, operatorID: operatorID}
	}
	return prepMsg, nil
}

func (msgContainer *MsgContainer) GetRound2Msg(operatorID uint32) (*Round2Message, error) {
	msg, err := msgContainer.GetMessage(Round2, operatorID)
	if err != nil {
		return nil, err
	}
	prepMsg, _ := msg.(*Round2Message)
	if prepMsg == nil {
		return nil, ErrMsgNil{round: Round2, operatorID: operatorID}
	}
	return prepMsg, nil
}

func (msgContainer *MsgContainer) GetBlameMsg(operatorID uint32) (*BlameMessage, error) {
	msg, err := msgContainer.GetMessage(Blame, operatorID)
	if err != nil {
		return nil, err
	}
	prepMsg, _ := msg.(*BlameMessage)
	if prepMsg == nil {
		return nil, ErrMsgNil{round: Blame, operatorID: operatorID}
	}
	return prepMsg, nil
}

func (msgContainer *MsgContainer) GetMessage(round ProtocolRound, operatorID uint32) (interface{}, error) {
	msg, ok := msgContainer.msgs[round][operatorID]
	if !ok {
		return nil, ErrMsgNotFound{round: round, operatorID: operatorID}
	}
	pm := &ProtocolMsg{}
	if err := pm.Decode(msg.Message.Data); err != nil {
		return nil, err
	}
	switch round {
	case Preparation:
		return pm.PreparationMessage, nil
	case Round1:
		return pm.Round1Message, nil
	case Round2:
		return pm.Round2Message, nil
	case Blame:
		return pm.BlameMessage, nil
	default:
		return nil, dkg.ErrInvalidRound{}
	}
}

func (msgContainer *MsgContainer) AllMessagesForRound(round ProtocolRound) map[uint32]*dkg.SignedMessage {
	return msgContainer.msgs[round]
}

func (msgContainer *MsgContainer) allMessagesReceivedFor(round ProtocolRound, operators []uint32) bool {
	for _, operatorID := range operators {
		if _, ok := msgContainer.msgs[round][operatorID]; !ok {
			return false
		}
	}
	return true
}

type ErrMsgNotFound struct {
	round      ProtocolRound
	operatorID uint32
}

func (e ErrMsgNotFound) Error() string {
	return fmt.Sprintf("message for operatorID %d and round %d not found\n", e.operatorID, e.round)
}

type ErrMsgNil struct {
	round      ProtocolRound
	operatorID uint32
}

func (e ErrMsgNil) Error() string {
	return fmt.Sprintf("message for operatorID %d and round %d is nil\n", e.operatorID, e.round)
}
