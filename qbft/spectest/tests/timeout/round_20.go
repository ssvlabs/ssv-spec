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
		PostRoot: "d676b1182d7b6a8fabcdab524c90fceb308fccf55aa0215925802132ac02ee25",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    21,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    21,
		},
	}
}
