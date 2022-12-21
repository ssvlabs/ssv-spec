package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Valid tests a valid start instance
func Valid() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance valid",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "7b74be21fcdae2e7ed495882d1a499642c15a7f732f210ee84fb40cc97d1ce96",
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
