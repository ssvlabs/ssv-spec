package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// InvalidPrepareJustificationValue tests a proposal for > 1 round, prepared previously but one of the prepare justifications has value != highest prepared
func InvalidPrepareJustificationValue() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := invalidPrepareJustificationValueStateComparison()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessageWithParams(
			// TODO: different value instead of wrong root
			ks.OperatorKeys[3], types.OperatorID(3), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingIdentifier, testingutils.DifferentRoot),
	}
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[3], types.OperatorID(3), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "invalid prepare justification value",
		Pre:            pre,
		PostRoot:       sc.Root(),
		PostState:      sc.ExpectedState,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: change round msg not valid: round change justification invalid: proposed data mismatch",
	}
}

func invalidPrepareJustificationValueStateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State

	return &comparable.StateComparison{ExpectedState: state}
}
