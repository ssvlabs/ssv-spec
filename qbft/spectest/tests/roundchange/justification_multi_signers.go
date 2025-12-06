package roundchange

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// JustificationMultiSigners tests a single prepare justification msg with multiple signers
func JustificationMultiSigners() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMultiSignerMessage(
			[]*rsa.PrivateKey{ks.OperatorKeys[2], ks.OperatorKeys[3]},
			[]types.OperatorID{2, 3},
		),
	}
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	test := tests.NewMsgProcessingSpecTest(
		"justification multi signer",
		testdoc.RoundChangeJustificationMultiSignersDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		types.MessageAllowsOneSignerOnlyErrorCode,
		nil,
		ks,
	)

	return test
}
