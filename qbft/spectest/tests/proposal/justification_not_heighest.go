package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// JustificationsNotHeighest tests a proposal for > 1 round, prepared previously with rc justification prepares at different heights but the prepare justification is not the highest
func JustificationsNotHeighest() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 3
	ks := testingutils.Testing4SharesSet()

	prepareMsgs1 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}
	prepareMsgs2 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithParams(
			ks.Shares[1], types.OperatorID(1), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs1)),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.Shares[2], types.OperatorID(2), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			2, testingutils.MarshalJustifications(prepareMsgs2)),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 3),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(
			ks.Shares[1], types.OperatorID(1), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal justification not highest",
		Pre:            pre,
		PostRoot:       "2bfbfa0fe458080526c4c8bb90494f22609830d1b3a52d3439849ceb543902f4",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: signed prepare not valid",
	}
}
