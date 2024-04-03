package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongHeight tests a commit msg received with the wrong height
func WrongHeight() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], 3),

		testingutils.TestingCommitMessageWrongHeight(ks.NetworkKeys[1], 1),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong commit height",
		Pre:           pre,
		PostRoot:      "976cd5cecd58bba892a38ec0ef02b3aed4656fb89fef473d8af78fedf095439d",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
			testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		},
	}
}
