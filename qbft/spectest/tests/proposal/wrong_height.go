package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongHeight tests a proposal msg received with the wrong height
func WrongHeight() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithHeight(ks.OperatorKeys[1], types.OperatorID(1), 2),
	}

	test := tests.NewMsgProcessingSpecTest(
		"wrong proposal height",
		testdoc.ProposalWrongHeightDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: wrong msg height",
		nil,
	)

	test.SetPrivateKeys(ks)

	return test
}
