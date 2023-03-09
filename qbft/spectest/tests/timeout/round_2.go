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
		PostRoot: "beaef03728ef5dadfd5daf11046923930e787f0d77b824326bd7ec65c2338b45",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    3,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    3,
		},
	}
}
