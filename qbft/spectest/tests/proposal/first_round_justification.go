package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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

	test := tests.NewMsgProcessingSpecTest(
		"proposal first round justification",
		testdoc.ProposalFirstRoundJustificationDoc,
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

	test.SetPrivateKeys(ks)

	return test
}
