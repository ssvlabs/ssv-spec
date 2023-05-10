package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// Round20 tests calling UponRoundTimeout for round 20, testing state and broadcasted msgs
func Round20() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := round20StateComparison()

	pre := testingutils.BaseInstance()
	pre.State.Round = 20
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 20)

	return &SpecTest{
		Name:      "round 20",
		Pre:       pre,
		PostRoot:  sc.Root(),
		PostState: sc.ExpectedState,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    21,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    21,
		},
	}
}

func round20StateComparison() *comparable.StateComparison {
	state := testingutils.BaseInstance().State
	state.Round = 21

	return &comparable.StateComparison{ExpectedState: state}
}
