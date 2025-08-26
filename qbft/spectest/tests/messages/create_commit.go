package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateCommit tests creating a commit msg
func CreateCommit() tests.SpecTest {
	return tests.NewCreateMsgSpecTest(
		"create commit",
		testdoc.MessagesCreateCommitDoc,
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		10,
		nil,
		nil,
		tests.CreateCommit,
		"f48d140d8beee146db7decbad3dd17c99dba6f2a88069f4721818279dbfe380c",
		nil,
		"",
		nil,
	)
}
