package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NotPreparedPreviouslyJustification tests a proposal for > 1 round, not prepared previously
func NotPreparedPreviouslyJustification() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 5
	ks := testingutils.Testing4SharesSet()
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 5),
	}

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 5,
			testingutils.MarshalJustifications(rcMsgs)),
	}

	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 5),
	}

	return tests.NewMsgProcessingSpecTest(
		"proposal justification (not prepared)",
		testdoc.ProposalNotPreparedPreviouslyJustificationDoc,
		pre,
		"",
		nil,
		inputMessages,
		outputMessages,
		"",
		nil,
	)
}
