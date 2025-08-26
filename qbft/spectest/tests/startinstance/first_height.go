package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FirstHeight tests a starting the first instance
func FirstHeight() tests.SpecTest {
	return tests.NewControllerSpecTest(
		"start instance first height",
		testdoc.StartInstanceFirstHeightDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
		nil,
		"",
		nil,
		nil,
	)
}
