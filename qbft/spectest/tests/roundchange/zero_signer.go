package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ZeroSigner tests a signed round change msg with signer 0
func ZeroSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(0), 2),
	}

	return tests.NewMsgProcessingSpecTest(
		"round change zero signer",
		"Test round change message with signer ID 0, expecting validation error.",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: invalid SignedSSVMessage: signer ID 0 not allowed",
		nil,
	)
}
