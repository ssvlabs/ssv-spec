package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// CreateRoundChangePreviouslyPrepared tests creating a round change msg,previously prepared
func CreateRoundChangePreviouslyPrepared() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create round change previously prepared",
		Value:      []byte{1, 2, 3, 4},
		PrepareJustifications: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
			}),
		},
		ExpectedRoot: "de421aca0c42404cea60b6a9dea40457831eed611b78dd529e022fc54e81af44",
	}
}
