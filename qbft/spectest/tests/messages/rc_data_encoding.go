package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataEncoding tests encoding RoundChangeData
func RoundChangeDataEncoding() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	signMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: nil},
	})
	signMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: nil},
	})
	signMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: nil},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         qbft.FirstRound,
		Input:         &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: 2,
	})

	rcMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		signMsg,
		signMsg2,
		signMsg3,
	}

	r, _ := rcMsg.GetRoot()
	b, _ := rcMsg.Encode()

	return &tests.MsgSpecTest{
		Name: "round change data encoding",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
				Data: b,
			},
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][]byte{
			r,
		},
	}
}
