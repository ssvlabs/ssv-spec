package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreatePrepare tests creating a prepare msg
func CreatePrepare() tests.SpecTest {
	return tests.NewCreateMsgSpecTest(
		"create prepare",
		testdoc.MessagesCreatePrepareDoc,
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		10,
		nil,
		nil,
		tests.CreatePrepare,
		"edf017cce3ba879ef8c16bda28c30c776ae6c522b41c7d0eb1cb1cf29f719613",
		nil,
		"",
		nil,
	)
}
