package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// QuorumOrder1 tests LIVENESS where the rc quorum msgs in different order
func QuorumOrder1() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	sc := quorumOrder1StateComparison()

	ks := testingutils.Testing4SharesSet()
	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change quorum order 1",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
				testingutils.TestingQBFTRootData,
				testingutils.MarshalJustifications(msgs), testingutils.MarshalJustifications(prepareMsgs)),
		},
	}
}

func quorumOrder1StateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
	}

	instance := &qbft.Instance{
		State: &qbft.State{
			Share: testingutils.TestingOperator(testingutils.Testing4SharesSet()),
			ID:    testingutils.TestingIdentifier,
			Round: 2,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	return &comparable.StateComparison{ExpectedState: instance.State}
}
