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
				ControllerPostRoot: "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
