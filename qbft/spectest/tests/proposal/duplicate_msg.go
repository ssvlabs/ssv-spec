package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate proposal msg processing
func DuplicateMsg() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	}

	return tests.NewMsgProcessingSpecTest(
		"proposal duplicate message",
		"Test processing of duplicate proposal messages from the same signer",
		pre,
		"",
		nil,
		inputMessages,
		outputMessages,
		"invalid signed message: proposal is not valid with current state",
		nil,
	)
}
