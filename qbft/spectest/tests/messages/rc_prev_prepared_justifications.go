package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RoundChangePrePreparedJustifications tests valid justified change round (prev prepared)
func RoundChangePrePreparedJustifications() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}

	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.OperatorKeys[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs))

	test := tests.NewMsgSpecTest(
		"rc prev prepared justifications",
		testdoc.MessagesRCPrevPreparedJustificationsDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		0,
		ks,
	)

	return test
}
