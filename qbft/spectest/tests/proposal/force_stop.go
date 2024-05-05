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
		PostRoot:      "eaa7264b5d6f05cfcdec3158fcae4ff58c3de1e7e9e12bd876177a58686994d4",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
