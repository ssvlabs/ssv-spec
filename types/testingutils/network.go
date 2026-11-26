package testingutils

import (
	"crypto/rsa"
	"fmt"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SignedSSVMessage
	OperatorID      types.OperatorID
	OperatorSK      *rsa.PrivateKey
	// failPartialSigBroadcast holds the partial-signature types whose containers fail to broadcast.
	failPartialSigBroadcast []types.PartialSigMsgType
}

func NewTestingNetwork(operatorID types.OperatorID, sk *rsa.PrivateKey) *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SignedSSVMessage, 0),
		OperatorID:      operatorID,
		OperatorSK:      sk,
	}
}

// SetFailPartialSigBroadcast FOR TESTING ONLY!! sets the partial-signature types whose containers fail to
// broadcast
func (net *TestingNetwork) SetFailPartialSigBroadcast(msgTypes []types.PartialSigMsgType) {
	net.failPartialSigBroadcast = msgTypes
}

// Broadcast records the message, or fails without recording it under SetFailPartialSigBroadcast.
func (net *TestingNetwork) Broadcast(msgID types.MessageID, message *types.SignedSSVMessage) error {
	if len(net.failPartialSigBroadcast) > 0 && message.SSVMessage.MsgType == types.SSVPartialSignatureMsgType {
		msgs := &types.PartialSignatureMessages{}
		if err := msgs.Decode(message.SSVMessage.Data); err != nil {
			return err
		}
		for _, msgType := range net.failPartialSigBroadcast {
			if msgs.Type == msgType {
				return fmt.Errorf("forced broadcast failure")
			}
		}
	}
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}

func ConvertBroadcastedMessagesToSSVMessages(signedMessages []*types.SignedSSVMessage) []*types.SSVMessage {
	ret := make([]*types.SSVMessage, 0)
	for _, msg := range signedMessages {
		ret = append(ret, msg.SSVMessage)
	}
	return ret
}

func (net *TestingNetwork) ExtractProposalMessages() []*types.SignedSSVMessage {
	ret := make([]*types.SignedSSVMessage, 0)
	for _, msg := range net.BroadcastedMsgs {
		if msg.SSVMessage.MsgType == types.SSVConsensusMsgType {
			qbftMsg := &qbft.Message{}
			if err := qbftMsg.Decode(msg.SSVMessage.Data); err == nil {
				if qbftMsg.MsgType == qbft.ProposalMsgType {
					ret = append(ret, msg)
				}
			}
		}
	}
	return ret
}

func (net *TestingNetwork) GetFullDataFromBroadcastedMessages() [][]byte {
	return GetFullDataFromMessages(net.ExtractProposalMessages())
}

func GetFullDataFromMessages(msgs []*types.SignedSSVMessage) [][]byte {
	ret := make([][]byte, 0)
	for _, msg := range msgs {
		ret = append(ret, msg.FullData)
	}
	return ret
}
