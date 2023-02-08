package startinstance

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Valid tests a valid start instance
func Valid() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance valid",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
