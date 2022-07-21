package testutils

import (
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
	"sync"
)

type MockNetwork struct {
	Broadcasted []*dkgtypes.Message
	mutex       sync.Mutex
}

func newMockNetwork() *MockNetwork {
	return &MockNetwork{
		Broadcasted: make([]*dkgtypes.Message, 0),
		mutex:       sync.Mutex{},
	}
}

func (m *MockNetwork) StreamDKGOutput(output map[types.OperatorID]*dkgtypes.SignedOutput) error {
	panic("implement me")
}

func (m *MockNetwork) Broadcast(msg types.Encoder) error {
	dkgMsg:=msg.(*dkgtypes.Message)
	m.Broadcasted = append(m.Broadcasted, dkgMsg)
	return nil
}
