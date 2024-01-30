package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// InvalidProposalJustificationValue tests a proposal for > 1 round, prepared previously but one of the prepare justifications has value != highest prepared
func InvalidProposalJustificationValue() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := invalidProposalJustificationValueStateComparison()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessageWithParams(
			// TODO: different value instead of wrong root
			ks.Shares[3], types.OperatorID(3), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingIdentifier, testingutils.DifferentRoot),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[3], types.OperatorID(3), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
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
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: No highest prepared round-change matches prepared messages",
	}
}

func invalidProposalJustificationValueStateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State

	return &comparable.StateComparison{ExpectedState: state}
}
