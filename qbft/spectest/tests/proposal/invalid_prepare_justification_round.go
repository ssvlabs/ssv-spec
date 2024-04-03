package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// InvalidPrepareJustificationRound tests a proposal for > 1 round, prepared previously but one of the prepare justifications has round != highest prepared round
func InvalidPrepareJustificationRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := invalidPrepareJustificationRoundStateComparison()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),
	}
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.NetworkKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.NetworkKeys[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.NetworkKeys[3], types.OperatorID(3), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.NetworkKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "invalid prepare justification round",
		Pre:            pre,
		PostRoot:       sc.Root(),
		PostState:      sc.ExpectedState,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: change round msg not valid: round change justification invalid: wrong msg round",
	}
}

func invalidPrepareJustificationRoundStateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State

	return &comparable.StateComparison{ExpectedState: state}
}
