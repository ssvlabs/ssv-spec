package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// AllHighestPreparedRCWithDiffData tests a proposal for r > 1 with all highest prepared RC with data != than prepare messages data
func AllHighestPreparedRCWithDiffData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithFullData(ks.Shares[1], types.OperatorID(1), testingutils.DifferentFullData),
		testingutils.TestingPrepareMessageWithFullData(ks.Shares[2], types.OperatorID(2), testingutils.DifferentFullData),
		testingutils.TestingPrepareMessageWithFullData(ks.Shares[3], types.OperatorID(3), testingutils.DifferentFullData),
	}

	// Single prepared RC
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], 1, 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, 1, testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[2], 2, 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, 1, testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[3], 3, 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, 1, testingutils.MarshalJustifications(prepareMsgs)),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "all highest prepared rc with diff data",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: No highest prepared round-change matches prepared messages",
	}
}
