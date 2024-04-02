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
		PostRoot:      "35d96d47901971b46055a010e136a97c71644d7d8bee368a8c1f29c149d3fa6c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: did not receive proposal for this round",
	}
}
