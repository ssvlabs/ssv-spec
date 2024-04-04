package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoPrevAcceptedProposal tests a commit msg received without a previous accepted proposal
func NoPrevAcceptedProposal() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = nil
	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessage(ks.Shares[1], 1),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "no previous accepted proposal",
		Pre:           pre,
		PostRoot:      "eaa7264b5d6f05cfcdec3158fcae4ff58c3de1e7e9e12bd876177a58686994d4",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: did not receive proposal for this round",
	}
}
