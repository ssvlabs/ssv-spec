package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongData2 tests a single commit received with a different commit data than the prepared data
func WrongData2() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingCommitMessageWrongRoot(ks.Shares[1], 1),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "commit data != prepared data",
		Pre:           pre,
		PostRoot:      "fd088c77032874833a90280172966c4a24cc0e661bb132ab028ff3b2b36a99f2",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
		ExpectedError: "invalid signed message: proposed data mistmatch",
	}
}
