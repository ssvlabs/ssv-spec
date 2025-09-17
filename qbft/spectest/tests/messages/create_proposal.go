package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateProposal tests creating a proposal msg, not previously prepared
func CreateProposal() tests.SpecTest {
	return tests.NewCreateMsgSpecTest(
		"create proposal",
		testdoc.MessagesCreateProposalDoc,
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		10,
		nil,
		nil,
		tests.CreateProposal,
		"a3225c24c7a759c09ecc69d4df1d5727a9ea417685b0084809377a285693cd48",
		nil,
		"",
		nil,
	)
}
