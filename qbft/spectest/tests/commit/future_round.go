package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureRound tests a commit msg received with a future round, should error
func FutureRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1)

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWithRound(ks.OperatorKeys[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit future round",
		Pre:           pre,
		PostRoot:      "fc68d10ca3bc1250da3aa789d4feae5fa197d08964c26c9f16cdfd39931d38ce",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg round",
	}
}
