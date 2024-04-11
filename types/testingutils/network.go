package testingutils

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SignedSSVMessage
	OperatorID      types.OperatorID
	OperatorSK      *rsa.PrivateKey
}

func NewTestingNetwork(operatorID types.OperatorID, sk *rsa.PrivateKey) *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SignedSSVMessage, 0),
		OperatorID:      operatorID,
		OperatorSK:      sk,
	}
}

func (net *TestingNetwork) Broadcast(message *types.SignedSSVMessage) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}

func ConvertBroadcastedMessagesToSSVMessages(signedMessages []*types.SignedSSVMessage) []*types.SSVMessage {
	ret := make([]*types.SSVMessage, 0)
	for _, msg := range signedMessages {
		ssvMsg := &types.SSVMessage{}
		err := ssvMsg.Decode(msg.Data)
		if err != nil {
			panic(err)
		}
		ret = append(ret, ssvMsg)
	}
	return ret
}
