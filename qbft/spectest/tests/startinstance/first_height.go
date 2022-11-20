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
				InputValue:         []byte{1, 2, 3, 4},
				DecidedVal:         nil,
				ControllerPostRoot: "3cf8cd7d050943781b25a611cc7b51bc051608f316556a843936b62c276984cb",
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
