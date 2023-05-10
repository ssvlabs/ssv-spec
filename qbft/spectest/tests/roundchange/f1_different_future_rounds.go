package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// F1DifferentFutureRounds tests f+1 speedup with one rc prev prepared
func F1DifferentFutureRounds() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := f1DifferentFutureRoundsStateComparison()

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
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
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

func f1DifferentFutureRoundsStateComparison() *comparable.StateComparison {
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

	instance := &qbft.Instance{
		State: &qbft.State{
			Share: testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:    testingutils.TestingIdentifier,
			Round: 5,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	return &comparable.StateComparison{ExpectedState: instance.State}
}
