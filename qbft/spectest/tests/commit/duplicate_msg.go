package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate commit msg processing
func DuplicateMsg() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], 1)

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMessage(ks.Shares[1], 1),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate commit message",
		Pre:           pre,
		PostRoot:      "f1ec378afe77b918d876515e14b9e4cbd911a17a8a1c584e613afc7e32245fc7",
		InputMessages: msgs,
	}
}
