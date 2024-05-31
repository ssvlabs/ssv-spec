package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// FutureRound tests a commit msg received with a future round, should error
func FutureRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1)

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessageWithRound(ks.OperatorKeys[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit future round",
		Pre:           pre,
		PostRoot:      "d3c3540b20d61319771ef7b01c26d172fbabeabad5c09b1415c72a34e0fd145f",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg round",
	}
}
