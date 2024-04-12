package commit

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
		PostRoot:      "167c1835a17bab210547283205e8e9cc754cb0c8a7fcdfcee57a63315ff63378",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
