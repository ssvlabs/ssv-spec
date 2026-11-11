package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// JustificationNestedPrepareJustification tests a prepared round change whose justification includes a prepare
// message that itself carries prepare justifications
func JustificationNestedPrepareJustification() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.SignQBFTMsg(ks.OperatorKeys[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.TestingIdentifier,
			Root:       testingutils.TestingQBFTRootData,
			PrepareJustification: testingutils.MarshalJustifications([]*types.SignedSSVMessage{
				testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			}),
		}),
	}
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return tests.NewMsgProcessingSpecTest(
		"justification nested prepare justification",
		testdoc.RoundChangeJustificationNestedPrepareJustificationDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		types.UnexpectedJustificationErrorCode,
		nil,
		ks,
	)
}
