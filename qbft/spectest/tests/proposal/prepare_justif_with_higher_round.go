package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PrepareJustifWithHigherRoundThanHighestPreparedRC tests a proposal for r > 1 with prepare messages with a higher round than the highest prepared RC
func PrepareJustifWithHigherRoundThanHighestPreparedRC() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	prepareMsgsForRound2 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}

	prepareMsgsForRound2Encoded, err := qbft.MarshalJustifications(prepareMsgsForRound2)
	if err != nil {
		panic(err.Error())
	}

	prepareMsgsForRound3 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 3),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 3),
	}

	// No prepared RC
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], 1, 4, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			2, prepareMsgsForRound2Encoded),
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[2], 2, 4, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			2, prepareMsgsForRound2Encoded),
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[3], 3, 4, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			2, prepareMsgsForRound2Encoded),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 4, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgsForRound3),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "prepare justification with higher round than highest prepared RC",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: No highest prepared round-change matches prepared messages",
	}
}
