package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// Round15 tests calling UponRoundTimeout for round 15, testing state and broadcasted msgs
func Round15() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := round20StateComparison()

	pre := testingutils.BaseInstance()
	pre.State.Round = 20
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 20)

	return &SpecTest{
		Name:      "round 15",
		Pre:       pre,
		PostRoot:  sc.Root(),
		PostState: sc.ExpectedState,

		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 0,
			Round:    0,
		},
		ExpectedError: "instance stopped processing timeouts",
	}
}

func round20StateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	state := testingutils.BaseInstance().State
	state.Round = 20
	state.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 20)

	return &comparable.StateComparison{ExpectedState: state}
}
