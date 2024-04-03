package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoPrepareQuorum tests a commit msg received without a previous prepare quorum
func NoPrepareQuorum() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], 2),
		// only 2 out of 4 prepare msgs

		testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingCommitMessage(ks.NetworkKeys[2], 2),
		testingutils.TestingCommitMessage(ks.NetworkKeys[3], 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "commit no prepare quorum",
		Pre:           pre,
		PostRoot:      "254e9491d5e3167c239ec04149e7b6d29402a6f8638010c35e04f1cea7f8f7e0",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
		},
	}
}
