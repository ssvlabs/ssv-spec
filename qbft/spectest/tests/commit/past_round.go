package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRound tests a commit msg with past round, should process but not decide
func PastRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.NetworkKeys[1], 1, 5)
	pre.State.Round = 5

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit past round",
		Pre:           pre,
		PostRoot:      "929e208f33ae4f872976a962256ddddebfb625905ff02903315586688ec1562d",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
