package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1DifferentFutureRoundsNotPrepared tests f+1 speedup (not prev prepared)
func F1DifferentFutureRoundsNotPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := f1DifferentFutureRoundsNotPreparedStateComparison()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 10),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change f+1 not prepared",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], types.OperatorID(1), 5, qbft.FirstHeight,
				[32]byte{}, 0, [][]byte{}),
		},
	}
}

func f1DifferentFutureRoundsNotPreparedStateComparison() *qbftcomparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	state := testingutils.BaseInstance().State
	state.Round = 5

	state.RoundChangeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.Round(5): {
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 5),
		},
		qbft.Round(10): {
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 10),
		},
	}}

	return &qbftcomparable.StateComparison{ExpectedState: state}
}
