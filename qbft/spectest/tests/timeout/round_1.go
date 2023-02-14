package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Round1 tests calling UponRoundTimeout for round 1, testing state and broadcasted msgs
func Round1() *SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	return &SpecTest{
		Name:     "round 1",
		Pre:      pre,
		PostRoot: "f807a3d5343ca20e3b757fa63eae9b9dd70e09e03249048badebf54e62290103",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    2,
		},
	}
}
