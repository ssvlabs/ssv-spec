package prepare

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongHeight tests prepare msg received with the wrong height
func WrongHeight() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithHeight(ks.OperatorKeys[1], types.OperatorID(1), 2),
	}

	test := tests.NewMsgProcessingSpecTest(
		"prepare wrong height",
		testdoc.PrepareWrongHeightDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		types.WrongMessageHeightErrorCode,
		nil,
		ks,
	)

	return test
}
