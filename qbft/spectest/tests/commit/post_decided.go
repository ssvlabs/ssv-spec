package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PostDecided tests processing a commit msg after instance decided
func PostDecided() *tests.MsgProcessingSpecTest {
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
		PostRoot:       "76a8ded050dc7578ec51f15e0f1e5e7dd63451d0b688bc4f3a818b794008a8a5",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}
