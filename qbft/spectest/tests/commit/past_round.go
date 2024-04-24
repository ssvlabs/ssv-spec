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
		PostRoot:      "a8d66621100a3f07eef3412e1bebece90917c1d7116b95f857c070dc33b2e5ad",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
