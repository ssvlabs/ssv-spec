package timeout

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Round2 tests calling UponRoundTimeout for round 2, testing state and broadcasted msgs
func Round2() *SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      2,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
	})

	return &SpecTest{
		Name:     "round 2",
		Pre:      pre,
		PostRoot: "d76f5d27ebdc1f33ed4af370fc7edb8a117d29a759597dbdf45560095a28151e",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.RoundChangeMsgType,
				Height:     qbft.FirstHeight,
				Round:      3,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
			}),
		},
		ExpectedTimerState: &testingutils.TimerState{
			Timeouts: 1,
			Round:    3,
		},
	}
}
