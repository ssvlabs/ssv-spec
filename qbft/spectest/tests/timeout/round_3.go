package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Round3 tests calling UponRoundTimeout for round 3, testing state and broadcasted msgs
func Round3() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := round3StateComparison()

	pre := testingutils.BaseInstance()
	pre.State.Round = 3
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 3)

	return &SpecTest{
		Name:      "round 3",
		Pre:       pre,
		PostRoot:  sc.Root(),
		PostState: sc.ExpectedState,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    4,
				Identifier:               testingutils.TestingIdentifier,
				Root:                     [32]byte{},
				RoundChangeJustification: [][]byte{},
				PrepareJustification:     [][]byte{},
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    4,
		},
	}
}

func round3StateComparison() *qbftcomparable.StateComparison {
	state := testingutils.BaseInstance().State
	state.Round = 4

	return &qbftcomparable.StateComparison{ExpectedState: state}
}
