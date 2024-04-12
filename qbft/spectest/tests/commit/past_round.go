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

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], 1, 5)
	pre.State.Round = 5

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWithRound(ks.OperatorKeys[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit past round",
		Pre:           pre,
		PostRoot:      "b38b41f38999d8bc051bc3707297b57d70fdc916248d7bdd60c694f18f3ac0a1",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
