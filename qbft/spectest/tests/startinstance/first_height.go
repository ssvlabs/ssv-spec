package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FirstHeight tests a starting the first instance
func FirstHeight() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance first height",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
