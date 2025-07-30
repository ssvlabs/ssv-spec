package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ForceStop tests processing a proposal msg when instance force stopped
func ForceStop() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.ForceStop()

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	test := tests.NewMsgProcessingSpecTest(
		"force stop proposal message",
		testdoc.ProposalForceStopDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"instance stopped processing messages",
		nil,
		ks,
	)

	return test
}
