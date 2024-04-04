package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongHeight tests a commit msg received with the wrong height
func WrongHeight() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingCommitMessageWrongHeight(ks.Shares[1], 1),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong commit height",
		Pre:           pre,
		PostRoot:      "fd088c77032874833a90280172966c4a24cc0e661bb132ab028ff3b2b36a99f2",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
	}
}
