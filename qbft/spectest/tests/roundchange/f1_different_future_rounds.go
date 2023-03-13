package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1DifferentFutureRounds tests f+1 speedup with one rc prev prepared
func F1DifferentFutureRounds() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[2], types.OperatorID(2), 10,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change f+1 prepared",
		Pre:           pre,
		PostRoot:      "141d47f572e1620701feab85a1131206783f540e418693fabb80992a68ebaf16",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], types.OperatorID(1), 5, qbft.FirstHeight,
				[32]byte{}, 0, [][]byte{}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    qbft.Round(5),
		},
	}
}
