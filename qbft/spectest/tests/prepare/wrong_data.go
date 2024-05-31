package prepare

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
		PostRoot:      "141aa76b11c863caecab7640c39de47577964d1536706bd95fd4f8a5d6f4e059",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposed data mismatch",
	}
}
