package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongHeight tests a commit msg received with the wrong height
func WrongHeight() *tests.MsgProcessingSpecTest {
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
		PostRoot:      "d382895e5ab872f6ff0e2e06994864c04c19a9dc38559cbcc750ba6025178316",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
	}
}
