package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SingleHighestPreparedRCWithCorrectData tests a proposal for r > 1 with a single highest prepared RC with correct data
func SingleHighestPreparedRCWithCorrectData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithFullData(ks.Shares[1], types.OperatorID(1), testingutils.TestingQBFTFullData),
		testingutils.TestingPrepareMessageWithFullData(ks.Shares[2], types.OperatorID(2), testingutils.TestingQBFTFullData),
		testingutils.TestingPrepareMessageWithFullData(ks.Shares[3], types.OperatorID(3), testingutils.TestingQBFTFullData),
	}

	// Single prepared RC
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], 1, 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, 1, testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], 2, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], 3, 2),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "single highest prepared rc with correct data",
		Pre:           pre,
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		},
	}
}
