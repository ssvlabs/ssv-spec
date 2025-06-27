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
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return tests.NewMsgProcessingSpecTest(
		"proposal first round justification",
		"Test proposal justification for the first round, verifying that the proposer is correct and expecting prepare message broadcast.",
		pre,
		"",
		nil,
		inputMessages,
		outputMessages,
		"",
		&testingutils.TimerState{
			Timeouts: 0,
			Round:    qbft.NoRound,
		},
	)
}
