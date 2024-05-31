package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ForceStop tests processing a commit msg when instance force stopped
func ForceStop() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))
	pre.ForceStop()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "force stop commit message",
		Pre:           pre,
		PostRoot:      "d3c3540b20d61319771ef7b01c26d172fbabeabad5c09b1415c72a34e0fd145f",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
