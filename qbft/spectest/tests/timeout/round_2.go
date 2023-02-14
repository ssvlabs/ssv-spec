package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Round2 tests calling UponRoundTimeout for round 2, testing state and broadcasted msgs
func Round2() *SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 2)

	return &SpecTest{
		Name:     "round 2",
		Pre:      pre,
		PostRoot: "d76f5d27ebdc1f33ed4af370fc7edb8a117d29a759597dbdf45560095a28151e",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    3,
		},
	}
}
