package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1SpeedupPrevPrepared tests catching up to higher rounds via f+1 speedup, other peers are all at the same round (one prev prepared)
func F1SpeedupPrevPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()

	ks := testingutils.Testing4SharesSet()
	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], types.OperatorID(3)),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.NetworkKeys[3], types.OperatorID(3), 10,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "f+1 speed up prev prepared",
		Pre:           pre,
		PostRoot:      "8daa5ad4f87b49d0db9640b379736f063788e862f2ece35af707d2adf4a531d5",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingRoundChangeMessageWithParams(ks.NetworkKeys[1], types.OperatorID(1), 10, qbft.FirstHeight,
				[32]byte{}, 0, [][]byte{}),
		},
	}
}
