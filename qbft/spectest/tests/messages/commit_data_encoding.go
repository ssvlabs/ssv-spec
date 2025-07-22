package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CommitDataEncoding tests encoding CommitData
func CommitDataEncoding() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1))

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return tests.NewMsgSpecTest(
		"commit data nil or len 0",
		testdoc.MessagesCommitDataEncodingDoc,
		[]*types.SignedSSVMessage{msg},
		[][]byte{b},
		[][32]byte{r},
		"",
	)
}
