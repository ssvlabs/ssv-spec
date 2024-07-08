package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests a valid start instance
func Valid() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance valid",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
