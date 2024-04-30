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
		PostRoot:      "fc68d10ca3bc1250da3aa789d4feae5fa197d08964c26c9f16cdfd39931d38ce",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
	}
}
