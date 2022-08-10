package testingutils

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs    []*types.SSVMessage
	BroadcastedDKGMsgs []*dkg.SignedMessage // GLNOTE: Should we use two fields or should make a new TestingNetwork instance?
	Outputs            map[types.OperatorID]*dkg.SignedOutput
}

func NewTestingNetwork() *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs:    make([]*types.SSVMessage, 0),
		BroadcastedDKGMsgs: make([]*dkg.SignedMessage, 0),
		Outputs:            make(map[types.OperatorID]*dkg.SignedOutput, 0),
	}
}

func (net *TestingNetwork) Broadcast(message types.Encoder) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message.(*types.SSVMessage))
	return nil
}

func (net *TestingNetwork) BroadcastDecided(msg types.Encoder) error {
	return nil
}

// StreamDKGOutput will stream to any subscriber the result of the DKG
func (net *TestingNetwork) StreamDKGOutput(output map[types.OperatorID]*dkg.SignedOutput) error {
	for id, signedOutput := range output {
		net.Outputs[id] = signedOutput
	}

	return nil
}

// BroadcastDKGMessage will broadcast a msg to the dkg network
func (net *TestingNetwork) BroadcastDKGMessage(msg *dkg.SignedMessage) error {
	net.BroadcastedDKGMsgs = append(net.BroadcastedDKGMsgs, msg)
	return nil
}
