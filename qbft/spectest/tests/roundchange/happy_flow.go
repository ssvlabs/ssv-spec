package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs,
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),

		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),

		testingutils.TestingCommitMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingCommitMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingCommitMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	)
	return &tests.MsgProcessingSpecTest{
		Name:          "round change happy flow",
		Pre:           pre,
		PostRoot:      "4bc58c6a479a8ef3ad8a39442d10501c2e6ce5ad03a81f3f67623dce474487e0",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
				[32]byte{}, 0, nil),
			testingutils.TestingProposalMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
				testingutils.MarshalJustifications(rcMsgs)),
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
			testingutils.TestingCommitMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    qbft.Round(2),
		},
	}
}
