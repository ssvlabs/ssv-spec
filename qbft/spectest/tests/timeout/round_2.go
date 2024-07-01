package timeout

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// Round2 tests calling UponRoundTimeout for round 2, testing state and broadcasted msgs
func Round2() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := round2StateComparison()

	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2))

	return &SpecTest{
		Name:      "round 2",
		Pre:       pre,
		PostRoot:  sc.Root(),
		PostState: sc.ExpectedState,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    3,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    3,
		},
	}
}

func round2StateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State
	state.Round = 3

	return &comparable.StateComparison{ExpectedState: state}
}
