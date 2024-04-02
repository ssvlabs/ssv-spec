package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests processing prepare msg after instance decided
func PostDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMessage(ks.Shares[2], 2),
		testingutils.TestingCommitMessage(ks.Shares[3], 3),

		testingutils.TestingPrepareMessage(ks.Shares[4], 4),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare post decided",
		Pre:           pre,
		PostRoot:      "17ba08e8b52a27154d4fcdb25f6a19f6c0a0ce0c0d509578f3b615bb3aa4436c",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
	}
}
