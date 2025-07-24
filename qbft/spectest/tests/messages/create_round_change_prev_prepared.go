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
		[32]byte{1, 2, 3, 4},
		nil,
		qbft.FirstRound,
		nil,
		prepareJustifications,
		tests.CreateRoundChange,
		"a6ffc48674f1522fb90aa7bde2aa76cac54480cf366cdd4afcd7f8b4d548809a",
		nil,
		"",
		ks,
	)

	return test
}
