package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1)),
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs,
		testingutils.TestingProposalMessageWithParams(ks.NetworkKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),

		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),

		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),
	)
	return &tests.MsgProcessingSpecTest{
		Name:          "round change happy flow",
		Pre:           pre,
		PostRoot:      "fbd68c0d90d12a2f3994e0af51d9d9dae5afedbcfebedcdcfca6f0b6b3b4c0ad",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
			testingutils.TestingRoundChangeMessageWithParams(ks.NetworkKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
				[32]byte{}, 0, nil),
			testingutils.TestingProposalMessageWithRoundAndRC(ks.NetworkKeys[1], types.OperatorID(1), 2,
				testingutils.MarshalJustifications(rcMsgs)),
			testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
			testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    qbft.Round(2),
		},
	}
}
