package timeout

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// Round5 tests calling UponRoundTimeout for round 5, testing state and broadcasted msgs
func Round5() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := round5StateComparison()

	pre := testingutils.BaseInstance()
	pre.State.Round = 5
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 5))

	return &SpecTest{
		Name:      "round 5",
		Pre:       pre,
		PostRoot:  sc.Root(),
		PostState: sc.ExpectedState,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    6,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    6,
		},
	}
}

func round5StateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State
	state.Round = 6

	return &comparable.StateComparison{ExpectedState: state}
}
