package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PeerPreparedDifferentHeights tests a round change quorum where peers prepared on different heights
func PeerPreparedDifferentHeights() tests.SpecTest {
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
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 3),
		testingutils.TestingRoundChangeMessageWithParamsAndFullData(
			ks.OperatorKeys[2], types.OperatorID(2), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs1), testingutils.TestingQBFTFullData),
		testingutils.TestingRoundChangeMessageWithParamsAndFullData(ks.OperatorKeys[3], types.OperatorID(3), 3, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, 2, testingutils.MarshalJustifications(prepareMsgs2), testingutils.TestingQBFTFullData),
	}

	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 3, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(inputMessages), testingutils.MarshalJustifications(prepareMsgs2)),
	}

	test := tests.NewMsgProcessingSpecTest(
		"round change peer prepared different heights",
		testdoc.RoundChangePeerPreparedDifferentHeightsDoc,
		pre,
		"",
		nil,
		inputMessages,
		outputMessages,
		"",
		nil,
		ks,
	)

	return test
}
