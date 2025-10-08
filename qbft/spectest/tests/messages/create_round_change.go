package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateRoundChange tests creating a round change msg, not previously prepared
func CreateRoundChange() tests.SpecTest {
	return tests.NewCreateMsgSpecTest(
		"create round change",
		testdoc.MessagesCreateRoundChangeDoc,
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		qbft.FirstRound,
		nil,
		nil,
		tests.CreateRoundChange,
		"a6ffc48674f1522fb90aa7bde2aa76cac54480cf366cdd4afcd7f8b4d548809a",
		nil,
		0,
		nil,
	)
}
