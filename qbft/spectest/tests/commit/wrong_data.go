package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongData1 tests commit msg with data != acceptedProposalData.Data
func WrongData1() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], 1)

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessageWrongRoot(ks.Shares[1], 1),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit data != acceptedProposalData.Data",
		Pre:           pre,
		PostRoot:      "35f7aaa4f57445a48f189b8c8edf66a3d0e3d54fde910097b6946ca6fa4d73ab",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposed data mistmatch",
	}
}
