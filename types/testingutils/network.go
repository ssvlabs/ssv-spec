package testingutils

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs       []*types.SSVMessage
	DKGOutputs            map[types.OperatorID]*dkg.SignedOutput
	SyncHighestDecidedCnt int
}

func NewTestingNetwork() *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SSVMessage, 0),
		DKGOutputs:      make(map[types.OperatorID]*dkg.SignedOutput, 0),
	}
}

func (net *TestingNetwork) Broadcast(message types.Encoder) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message.(*types.SSVMessage))
	return nil
}

func (net *TestingNetwork) BroadcastDecided(message types.Encoder) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message.(*types.SSVMessage))
	return nil
}

// StreamDKGOutput will stream to any subscriber the result of the DKG
func (net *TestingNetwork) StreamDKGOutput(output map[types.OperatorID]*dkg.SignedOutput) error {
	for id, signedOutput := range output {
		net.DKGOutputs[id] = signedOutput
	}

	return nil
}

func (net *TestingNetwork) SyncHighestDecided(identifier []byte) error {
	net.SyncHighestDecidedCnt++
	return nil
}

func (net *TestingNetwork) SyncHighestRoundChange(identifier []byte, height qbft.Height) error {
	return nil
}

//func (net *TestingNetwork) SyncHighestDecided() error {
//	return nil
//}

// BroadcastDKGMessage will broadcast a msg to the dkg network
func (net *TestingNetwork) BroadcastDKGMessage(msg *dkg.SignedMessage) error {
	data, err := msg.Encode()
	if err != nil {
		return err
	}
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, &types.SSVMessage{
		MsgType: types.DKGMsgType,
		MsgID:   types.MessageID{}, // TODO: what should we use for the MsgID?
		Data:    data,
	})
	return nil
}
