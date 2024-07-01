package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongData1 tests commit msg with data != acceptedProposalData.Data
func WrongData1() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWrongRoot(ks.OperatorKeys[1], 1),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit data != acceptedProposalData.Data",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposed data mismatch",
	}
}
