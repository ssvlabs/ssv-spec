package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests processing a commit msg after instance decided
func PostDecided() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], 1)

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMessage(ks.Shares[2], 2),
		testingutils.TestingCommitMessage(ks.Shares[3], 3),
		testingutils.TestingCommitMessage(ks.Shares[4], 4),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "post decided",
		Pre:            pre,
		PostRoot:       "d001f57a5f55a407a83fdc633d5cd2dd30fec9423b928b18c89a3880f8a8da62",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}
