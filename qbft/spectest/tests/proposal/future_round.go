package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureRound tests a proposal for state.ProposalAcceptedForCurrentRound != nil && signedProposal.Message.Round > state.Round
func FutureRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1))
	pre.State.Round = qbft.FirstRound

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 8),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 8),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 8),
	}

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 10),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.NetworkKeys[1], types.OperatorID(1), 10, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal future round prev prepared",
		Pre:           pre,
		PostRoot:      "1d6b7cd0d77efabb5755072f2b502baa33759c2d5ce8df098294bfeec231a560",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 10),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    qbft.Round(10),
		},
	}
}
