package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// RoundChangeDataEncoding tests encoding RoundChangeData
func RoundChangeDataEncoding() *tests.MsgSpecTest {
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data: testingutils.RoundChangePreparedDataBytes(
			[]byte{1, 2, 3, 4},
			2,
			[]*qbft.SignedMessage{
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				}),
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				}),
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				}),
			},
		),
	})

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "round change data encoding",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][]byte{
			r,
		},
	}
}
