package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateRoundChangePreviouslyPrepared tests creating a round change msg,previously prepared
func CreateRoundChangePreviouslyPrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	prepareJustifications := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}
	test := tests.NewCreateMsgSpecTest(
		"create round change previously prepared",
		testdoc.MessagesCreateRoundChangePrevPreparedDoc,
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		qbft.FirstRound,
		nil,
		prepareJustifications,
		tests.CreateRoundChange,
		"56037aa401ffe4de6c79c737c6f2fb39bd6d36d4e5c71894fdaac0cda7ab4d1a",
		nil,
		"",
		ks,
	)

	return test
}
