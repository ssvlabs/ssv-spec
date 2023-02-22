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
		PostRoot: "da48cb7ea2776bbef227ff0da21527f51ff62239b411aa150cb72bb3129155ea",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    6,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    6,
		},
	}
}
