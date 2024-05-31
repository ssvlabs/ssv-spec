package prepare

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
		PostRoot:      "d3c3540b20d61319771ef7b01c26d172fbabeabad5c09b1415c72a34e0fd145f",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg height",
	}
}
