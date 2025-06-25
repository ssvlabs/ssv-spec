package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// F1Speedup tests catching up to higher rounds via f+1 speedup, other peers are all at the same round
func F1Speedup() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 10),
	}

	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 10, qbft.FirstHeight,
			[32]byte{}, 0, [][]byte{}),
	}

	return tests.NewMsgProcessingSpecTest(
		"f+1 speed up",
		pre,
		"",
		nil,
		inputMessages,
		outputMessages,
		"",
		&testingutils.TimerState{
			Timeouts: 1,
			Round:    qbft.Round(10),
		},
	)
}
