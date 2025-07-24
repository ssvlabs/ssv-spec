package timeout

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// Round3 tests calling UponRoundTimeout for round 3, testing state and broadcasted msgs
func Round3() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := round3StateComparison()

	pre := testingutils.BaseInstance()
	pre.State.Round = 3
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 3))

	outputMessages := []*types.SignedSSVMessage{
		testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
			MsgType:                  qbft.RoundChangeMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    4,
			Identifier:               testingutils.TestingIdentifier,
			Root:                     [32]byte{},
			RoundChangeJustification: [][]byte{},
			PrepareJustification:     [][]byte{},
		}),
	}

	test := NewSpecTest(
		"round 3",
		"Test UponRoundTimeout for round 3, checks state transition and broadcasted round change message.",
		pre,
		sc.Root(),
		sc.ExpectedState,
		outputMessages,
		&testingutils.TimerState{
			Timeouts: 1,
			Round:    4,
		},
		"",
	)

	test.SetPrivateKeys(ks)

	return test
}

func round3StateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State
	state.Round = 4

	return &comparable.StateComparison{ExpectedState: state}
}
