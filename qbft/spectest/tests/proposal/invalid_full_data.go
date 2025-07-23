package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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
		testdoc.ProposalInvalidFullDataDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: H(data) != root",
		nil,
	)
}
