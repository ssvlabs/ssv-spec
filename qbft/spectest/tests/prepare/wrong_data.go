package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongData tests prepare msg with data != acceptedProposalData.Data
func WrongData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))

	msgs := []*types.SignedSSVMessage{
		// TODO: different value instead of wrong root
		testingutils.TestingPrepareMessageWrongRoot(ks.OperatorKeys[1], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare wrong data",
		Pre:           pre,
		PostRoot:      "2253eea5735c33797cd1f1a1e3ced2cb8b16ee1c78ae1747e18041b67216d622",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposed data mistmatch",
	}
}
