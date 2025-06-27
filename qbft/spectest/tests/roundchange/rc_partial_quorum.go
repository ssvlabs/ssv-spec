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

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 3),
	}

	outMsg := testingutils.TestingRoundChangeMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
		[32]byte{}, 0, [][]byte{})
	outMsg.FullData = []byte{}

	outputMessages := []*types.SignedSSVMessage{outMsg}

	return tests.NewMsgProcessingSpecTest(
		"round change partial quorum",
		"Test round change with partial quorum, checks timer state and output message.",
		pre,
		sc.Root(),
		sc.ExpectedState,
		inputMessages,
		outputMessages,
		"",
		&testingutils.TimerState{
			Timeouts: 1,
			Round:    qbft.Round(2),
		},
	)
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
