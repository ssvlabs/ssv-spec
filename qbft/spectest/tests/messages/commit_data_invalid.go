package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CommitDataInvalid tests commit data len == 0
func CommitDataInvalid() *tests.MsgSpecTest {
	keySet := testingutils.Testing4SharesSet()
	msg := testingutils.TestingCommitMessageWithParams(
		keySet.Shares[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.WrongRoot,
	)

	return &tests.MsgSpecTest{
		Name: "commit data invalid",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "CommitData data is invalid",
	}
}
