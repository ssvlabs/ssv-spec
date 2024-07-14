package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// JustificationsNotHeighest tests a proposal for > 1 round, prepared previously with rc justification prepares at different heights but the prepare justification is not the highest
func JustificationsNotHeighest() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 3
	ks := testingutils.Testing4SharesSet()

	prepareMsgs1 := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}
	prepareMsgs2 := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
	}
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithParams(
			ks.OperatorKeys[1], types.OperatorID(1), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs1)),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.OperatorKeys[2], types.OperatorID(2), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			2, testingutils.MarshalJustifications(prepareMsgs2)),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 3),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(
			ks.OperatorKeys[1], types.OperatorID(1), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal justification not highest",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: signed prepare not valid",
	}
}
