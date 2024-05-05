package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PastRound tests a commit msg with past round, should process but not decide
func PastRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], 1, 5)
	pre.State.Round = 5

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWithRound(ks.OperatorKeys[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit past round",
		Pre:           pre,
		PostRoot:      "3783417654934abf09b3d676e4a2553884911fc845701dcb19d5555bf8e5c929",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
