package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ForceStop tests processing a proposal msg when instance force stopped
func ForceStop() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.ForceStop()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "force stop proposal message",
		Pre:           pre,
		PostRoot:      "01489f7af13579b66ce3da156d4d10208c85a10365380f04e7b8d82d0a9679ce",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
