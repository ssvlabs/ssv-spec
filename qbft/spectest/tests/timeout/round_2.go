package timeout

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
		PostRoot: "5a71daf1a4ee817826596858f76e56a1a85aedd85b9c4e65e73fc4c4667e65b0",
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.RoundChangeMsgType,
				Height:     qbft.FirstHeight,
				Round:      3,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
			}),
		},
	}
}
