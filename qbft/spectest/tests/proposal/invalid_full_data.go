package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidFullData tests signed proposal with an invalid full data field (H(full data) != root)
func InvalidFullData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))
	msg.FullData = nil

	inputMessages := []*types.SignedSSVMessage{msg}

	return tests.NewMsgProcessingSpecTest(
		"invalid full data",
		"Test proposal message with invalid full data field where hash of full data does not equal the root, expecting validation error.",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: H(data) != root",
		nil,
	)
}
