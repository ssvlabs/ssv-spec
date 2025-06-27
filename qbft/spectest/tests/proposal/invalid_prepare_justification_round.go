package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// InvalidPrepareJustificationRound tests a proposal for > 1 round, prepared previously but one of the prepare justifications has round != highest prepared round
func InvalidPrepareJustificationRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := invalidPrepareJustificationRoundStateComparison()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
	}
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[3], types.OperatorID(3), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return tests.NewMsgProcessingSpecTest(
		"invalid prepare justification round",
		"Test proposal for round > 1 that was prepared previously but contains prepare justification with round different from the highest prepared round, expecting validation error.",
		pre,
		sc.Root(),
		sc.ExpectedState,
		inputMessages,
		nil,
		"invalid signed message: proposal not justified: change round msg not valid: round change justification invalid: wrong msg round",
		nil,
	)
}

func invalidPrepareJustificationRoundStateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State

	return &comparable.StateComparison{ExpectedState: state}
}
