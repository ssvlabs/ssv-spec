package testutils

import (
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"sync"
)

type MockNetwork struct {
	Broadcasted []*dkgtypes.Message
	Output      map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage
	mutex       sync.Mutex
}

func newMockNetwork() *MockNetwork {
	return &MockNetwork{
		Broadcasted: make([]*dkgtypes.Message, 0),
		Output:      map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage{},
		mutex:       sync.Mutex{},
	}
}

func (m *MockNetwork) StreamDKGOutput(output map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage) error {
	for id, message := range output {
		m.Output[id] = message
	}
	return nil
}

func (m *MockNetwork) Broadcast(msg types.Encoder) error {
	dkgMsg := msg.(*dkgtypes.Message)
	m.Broadcasted = append(m.Broadcasted, dkgMsg)
	return nil
}
