package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*qbft.SignedMessage
}

func NewTestingNetwork() qbft.Network {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*qbft.SignedMessage, 0),
	}
}

func (net *TestingNetwork) Broadcast(message types.Encoder) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message.(*qbft.SignedMessage))
	return nil
}

func (net *TestingNetwork) BroadcastDecided(msg types.Encoder) error {
	return nil
}
