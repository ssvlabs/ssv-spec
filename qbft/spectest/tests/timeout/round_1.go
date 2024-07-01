package timeout

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// Round1 tests calling UponRoundTimeout for round 1, testing state and broadcasted msgs
func Round1() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := round1StateComparison()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))

	return &SpecTest{
		Name:      "round 1",
		Pre:       pre,
		PostRoot:  sc.Root(),
		PostState: sc.ExpectedState,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    2,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    2,
		},
	}
}

func round1StateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State
	state.Round = 2

	return &comparable.StateComparison{ExpectedState: state}
}
