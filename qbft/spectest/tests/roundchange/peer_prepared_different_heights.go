package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PeerPreparedDifferentHeights tests a round change quorum where peers prepared on different heights
func PeerPreparedDifferentHeights() *tests.MsgProcessingSpecTest {
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
	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.Shares[2], types.OperatorID(2), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs1)),
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[3], types.OperatorID(3), 3, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, 2, testingutils.MarshalJustifications(prepareMsgs2)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change peer prepared different heights",
		Pre:           pre,
		PostRoot:      "79f065a6e1e3b11e0cacc0e46c4253116a08c183b937801a2d355b9d5bbee969",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 3, qbft.FirstHeight,
				testingutils.TestingQBFTRootData,
				testingutils.MarshalJustifications(msgs), testingutils.MarshalJustifications(prepareMsgs2)),
		},
	}
}
