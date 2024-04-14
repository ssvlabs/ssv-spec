package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongHeight tests prepare msg received with the wrong height
func WrongHeight() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithHeight(ks.OperatorKeys[1], types.OperatorID(1), 2),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare wrong height",
		Pre:           pre,
		PostRoot:      "167c1835a17bab210547283205e8e9cc754cb0c8a7fcdfcee57a63315ff63378",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
	}
}
