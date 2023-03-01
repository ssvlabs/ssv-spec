package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FirstHeight tests a starting the first instance
func FirstHeight() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance first height",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
