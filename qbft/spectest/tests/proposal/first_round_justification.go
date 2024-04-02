package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FirstRoundJustification tests proposal justification for first round (proposer is correct check)
func FirstRoundJustification() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal first round justification",
		Pre:           pre,
		PostRoot:      "47bc12be6e7a498d07c3e585664eedf6f8634da5cb9140fb7fad5e507758d0d8",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 0,
			Round:    qbft.NoRound,
		},
	}
}
