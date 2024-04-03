package testingutils

import (
	"github.com/bloxapp/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SignedSSVMessage
}

func NewTestingNetwork() *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SignedSSVMessage, 0),
	}
}

func (net *TestingNetwork) Broadcast(message *types.SignedSSVMessage) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}
