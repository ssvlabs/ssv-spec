package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests processing a commit msg after instance decided
func PostDecided() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.NetworkKeys[1], 1)

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingCommitMessage(ks.NetworkKeys[2], 2),
		testingutils.TestingCommitMessage(ks.NetworkKeys[3], 3),
		testingutils.TestingCommitMessage(ks.NetworkKeys[4], 4),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "post decided",
		Pre:            pre,
		PostRoot:       "b778306f5782f6df3bc1aadd173f2e5dc41fcec00fee887a59c372bc94be91c2",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
