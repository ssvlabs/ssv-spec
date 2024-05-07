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

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1)

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWrongRoot(ks.OperatorKeys[1], 1),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit data != acceptedProposalData.Data",
		Pre:           pre,
		PostRoot:      "2253eea5735c33797cd1f1a1e3ced2cb8b16ee1c78ae1747e18041b67216d622",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposed data mistmatch",
	}
}
