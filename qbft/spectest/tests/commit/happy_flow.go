package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a quorum of commits received after prepare quorum
func HappyFlow() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], 3),

		testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingCommitMessage(ks.NetworkKeys[2], 2),
		testingutils.TestingCommitMessage(ks.NetworkKeys[3], 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "commit happy flow",
		Pre:           pre,
		PostRoot:      "a298f278a78362257e233e1db8eeb44c2bb9b45a55bb3555928d6723231ebcd2",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
			testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		},
	}
}
