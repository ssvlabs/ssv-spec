package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// F1SpeedupPrevPrepared tests catching up to higher rounds via f+1 speedup, other peers are all at the same round (one prev prepared)
func F1SpeedupPrevPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()

	ks := testingutils.Testing4SharesSet()
	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[3], types.OperatorID(3), 10,
			testingutils.MarshalJustifications(prepareMsgs)),
	}
	outputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 10, qbft.FirstHeight,
			[32]byte{}, 0, [][]byte{}),
	}

	test := tests.NewMsgProcessingSpecTest(
		"f+1 speed up prev prepared",
		testdoc.RoundChangeF1SpeedupPrevPreparedDoc,
		pre,
		"",
		nil,
		inputMessages,
		outputMessages,
		0,
		nil,
		ks,
	)

	return test
}
