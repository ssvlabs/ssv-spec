package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// PreparedPreviouslyNoPrepareJustificationQuorum tests a proposal for > 1 round, prepared previously but without quorum of prepared msgs justification
func PreparedPreviouslyNoPrepareJustificationQuorum() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := preparedPreviouslyNoPrepareJustificationQuorumStateComparison()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
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

	test := tests.NewMsgProcessingSpecTest(
		"no prepare quorum (prepared)",
		testdoc.ProposalPreparedPreviouslyNoPrepareJustificationQuorumDoc,
		pre,
		sc.Root(),
		sc.ExpectedState,
		inputMessages,
		nil,
		types.JustificationsNoQuorumInvalidErrorCode,
		nil,
		ks,
	)

	return test
}

func preparedPreviouslyNoPrepareJustificationQuorumStateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State

	return &comparable.StateComparison{ExpectedState: state}
}
