package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SSVMessage
}

func NewTestingNetwork() qbft.Network {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SSVMessage, 0),
	}
}

func (net *TestingNetwork) Broadcast(message types.Encoder) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message.(*types.SSVMessage))
	return nil
}

func (net *TestingNetwork) BroadcastDecided(msg types.Encoder) error {
	return nil
}
