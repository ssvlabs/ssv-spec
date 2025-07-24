package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongHeight tests a round change msg with wrong height
func WrongHeight() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.OperatorKeys[1], types.OperatorID(1), 2, 2),
	}

	test := tests.NewMsgProcessingSpecTest(
		"round change invalid height",
		testdoc.RoundChangeWrongHeightDoc,
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
