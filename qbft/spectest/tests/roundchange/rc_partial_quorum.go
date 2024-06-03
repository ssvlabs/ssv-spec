package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// RoundChangePartialQuorum tests a round change msgs with partial quorum
func RoundChangePartialQuorum() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := roundChangePartialQuorumStateComparison()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 3),
	}

	outMsg := testingutils.TestingRoundChangeMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
		[32]byte{}, 0, [][]byte{})
	outMsg.FullData = []byte{}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change partial quorum",
		Pre:            pre,
		PostRoot:       sc.Root(),
		PostState:      sc.ExpectedState,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{outMsg},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    qbft.Round(2),
		},
	}
}

func roundChangePartialQuorumStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 3),
	}

	instance := &qbft.Instance{
		State: &qbft.State{
			CommitteeMember: testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
			ID:              testingutils.TestingIdentifier,
			Round:           2,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	return &comparable.StateComparison{ExpectedState: instance.State}
}
