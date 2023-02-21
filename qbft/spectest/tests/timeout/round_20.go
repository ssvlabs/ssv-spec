package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Round20 tests calling UponRoundTimeout for round 20, testing state and broadcasted msgs
func Round20() *SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 20
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 20)

	return &SpecTest{
		Name:     "round 20",
		Pre:      pre,
		PostRoot: "742d15da204c3e2ec7f78b6234dca7ec71ff3e0e8605639b97c3226894892fbc",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 21),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    21,
		},
	}
}
