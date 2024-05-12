package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FirstRoundJustification tests proposal justification for first round (proposer is correct check)
func FirstRoundJustification() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal first round justification",
		Pre:           pre,
		PostRoot:      "a218ef33486e9bdebef968f145f69a3e4ec9c4d824842297d3f6510d72ed5003",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 0,
			Round:    qbft.NoRound,
		},
	}
}
