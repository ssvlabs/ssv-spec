package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRound tests a commit msg with past round, should process but not decide
func PastRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], 1, 5)
	pre.State.Round = 5

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessageWithRound(ks.Shares[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit past round",
		Pre:           pre,
		PostRoot:      "873356f21c47c49ba0edb44b3bcfe9c6ace306cb28a4fa6e5d00697e61141685",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
