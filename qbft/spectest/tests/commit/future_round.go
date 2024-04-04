package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureRound tests a commit msg received with a future round, should error
func FutureRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], 1)

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessageWithRound(ks.Shares[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit future round",
		Pre:           pre,
		PostRoot:      "e3194c84f99e73171890f32848497b619050587254bf2315ed757095ced37839",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg round",
	}
}
