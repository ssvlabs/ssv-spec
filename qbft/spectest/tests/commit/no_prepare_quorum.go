package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoPrepareQuorum tests a commit msg received without a previous prepare quorum
func NoPrepareQuorum() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		// only 2 out of 4 prepare msgs

		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMessage(ks.Shares[2], 2),
		testingutils.TestingCommitMessage(ks.Shares[3], 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "commit no prepare quorum",
		Pre:           pre,
		PostRoot:      "a789f5c8934ec150137c83c940cbfbb057f564fd579432778b688cfe64a02260",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		},
	}
}
