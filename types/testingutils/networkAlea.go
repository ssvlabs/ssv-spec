package testingutils

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type TestingNetworkAlea struct {
	BroadcastedMsgs           []*types.SSVMessage
	SyncHighestDecidedCnt     int
	SyncHighestChangeRoundCnt int
	DecidedByRange            [2]alea.Height
}

func NewTestingNetworkAlea() *TestingNetworkAlea {
	return &TestingNetworkAlea{
		BroadcastedMsgs: make([]*types.SSVMessage, 0),
	}
}

func (net *TestingNetworkAlea) Broadcast(message *types.SSVMessage) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}

func (net *TestingNetworkAlea) SyncHighestDecided(identifier types.MessageID) error {
	net.SyncHighestDecidedCnt++
	return nil
}

//func (net *TestingNetworkAlea) SyncHighestDecided() error {
//	return nil
//}

// SyncDecidedByRange will sync decided messages from-to (including them)
func (net *TestingNetworkAlea) SyncDecidedByRange(identifier types.MessageID, from, to alea.Height) {
	net.DecidedByRange = [2]alea.Height{from, to}
}
