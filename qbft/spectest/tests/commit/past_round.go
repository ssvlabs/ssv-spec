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
		PostRoot:      "87283f01015e2decc6be336ecdab69a11e06ec153e0b27a479bf572d1b113f0c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
