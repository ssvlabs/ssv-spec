package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidFullData tests the signed commit with a valid full data field (H(full data) == root)
func ValidFullData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingCommitMessageWithFullData(ks.Shares[1], 1, testingutils.TestingQBFTFullData),
		testingutils.TestingCommitMessageWithFullData(ks.Shares[2], 2, testingutils.TestingQBFTFullData),
		testingutils.TestingCommitMessageWithFullData(ks.Shares[3], 3, testingutils.TestingQBFTFullData),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "commit with valid full data",
		Pre:           pre,
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
	}
}
