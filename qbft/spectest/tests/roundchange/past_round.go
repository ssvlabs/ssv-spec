package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRound tests a round change msg with past round
func PastRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 12
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], types.OperatorID(3)),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 5),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.NetworkKeys[1], types.OperatorID(1), 10,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change past round",
		Pre:            pre,
		PostRoot:       "708abf25fd6ea204c0b177e1b28fe0a6a0da3fb69da87181f73399ddde1a30b0",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: past round",
	}
}
