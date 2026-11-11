package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RCJustificationNestedJustification tests a proposal for round > 1 justified by a quorum of valid round changes plus
// a not prepared round change that itself carries round change justifications
func RCJustificationNestedJustification() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		testingutils.TestingRoundChangeMessageWithParams(ks.OperatorKeys[4], types.OperatorID(4), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, qbft.NoRound,
			testingutils.MarshalJustifications([]*types.SignedSSVMessage{
				testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			})),
	}

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(rcMsgs)),
	}

	return tests.NewMsgProcessingSpecTest(
		"proposal rc justification nested justification",
		testdoc.ProposalRCJustificationNestedJustificationDoc,
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
