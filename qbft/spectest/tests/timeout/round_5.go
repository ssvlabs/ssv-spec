package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Round5 tests calling UponRoundTimeout for round 5, testing state and broadcasted msgs
func Round5() *SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 5
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 5)

	return &SpecTest{
		Name:     "round 5",
		Pre:      pre,
		PostRoot: "538d592e46ebabfef1d14131a105e0a3daa532f5a50c6316d5d0c56b26cbe6ff",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 6),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    6,
		},
	}
}
