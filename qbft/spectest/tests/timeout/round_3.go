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
		PostRoot: "71f78b5d9365cb66875cc81bc984a0d0a8fd9c683243a981de397adabd238f7f",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    4,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    4,
		},
	}
}
