package common

import (
	"fmt"
	"sync"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/pkg/errors"
)

type IMsgContainer interface {
	SaveMsg(round ProtocolRound, msg *dkg.SignedMessage) (existingMessage *dkg.SignedMessage, err error)
	GetSignedMsg(round ProtocolRound, operatorID uint32) (*dkg.SignedMessage, error)
	AllMessagesForRound(round ProtocolRound) map[uint32]*dkg.SignedMessage
	AllMessagesReceivedFor(round ProtocolRound, operators []uint32) bool
	AllMessagesReceivedUpto(round ProtocolRound, operators []uint32, threshold uint64) bool
}

type MsgContainer struct {
	mu   *sync.Mutex
	msgs map[ProtocolRound]map[uint32]*dkg.SignedMessage
}

func NewMsgContainer() IMsgContainer {
	m := make(map[ProtocolRound]map[uint32]*dkg.SignedMessage)
	for _, round := range rounds {
		m[round] = make(map[uint32]*dkg.SignedMessage)
	}
	return &MsgContainer{msgs: m, mu: new(sync.Mutex)}
}

func (msgContainer *MsgContainer) SaveMsg(round ProtocolRound, msg *dkg.SignedMessage) (existingMessage *dkg.SignedMessage, err error) {
	msgContainer.mu.Lock()
	defer msgContainer.mu.Unlock()

	existingMessage, exists := msgContainer.msgs[round][uint32(msg.Signer)]
	if exists {
		return existingMessage, errors.New("msg already exists")
	}
	msgContainer.msgs[round][uint32(msg.Signer)] = msg
	return nil, nil
}

func (msgContainer *MsgContainer) GetSignedMsg(round ProtocolRound, operatorID uint32) (*dkg.SignedMessage, error) {
	msgContainer.mu.Lock()
	defer msgContainer.mu.Unlock()

	signedMsg, exist := msgContainer.msgs[round][operatorID]
	if !exist {
		return nil, ErrMsgNotFound{Round: round, OperatorID: operatorID}
	}
	return signedMsg, nil
}

func (msgContainer *MsgContainer) AllMessagesForRound(round ProtocolRound) map[uint32]*dkg.SignedMessage {
	msgContainer.mu.Lock()
	defer msgContainer.mu.Unlock()

	return msgContainer.msgs[round]
}

func (msgContainer *MsgContainer) AllMessagesReceivedFor(round ProtocolRound, operators []uint32) bool {
	msgContainer.mu.Lock()
	defer msgContainer.mu.Unlock()

	for _, operatorID := range operators {
		if _, ok := msgContainer.msgs[round][uint32(operatorID)]; !ok {
			return false
		}
	}
	return true
}

func (msgContainer *MsgContainer) AllMessagesReceivedUpto(round ProtocolRound, operators []uint32, threshold uint64) bool {
	msgContainer.mu.Lock()
	defer msgContainer.mu.Unlock()

	totalMsgsRecieved := uint64(0)
	for _, operatorID := range operators {
		if _, ok := msgContainer.msgs[round][uint32(operatorID)]; ok {
			totalMsgsRecieved += 1
		}
	}
	return totalMsgsRecieved >= threshold
}

type ErrMsgNotFound struct {
	Round      ProtocolRound
	OperatorID uint32
}

func (e ErrMsgNotFound) Error() string {
	return fmt.Sprintf("message for operatorID %d and round %d not found\n", e.OperatorID, e.Round)
}

type ErrMsgNil struct {
	Round      ProtocolRound
	OperatorID uint32
}

func (e ErrMsgNil) Error() string {
	return fmt.Sprintf("message for operatorID %d and round %d is nil\n", e.OperatorID, e.Round)
}
