package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Round3 tests calling UponRoundTimeout for round 3, testing state and broadcasted msgs
func Round3() *SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 3
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 3)

	return &SpecTest{
		Name:     "round 3",
		Pre:      pre,
		PostRoot: "d3989251b49ba2ca86166038c3efc762e2e20a5467289d127223171f16f5eda3",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 4),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    4,
		},
	}
}
