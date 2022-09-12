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
		Input:  []byte{1, 2, 3, 4},
	})
	signMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  []byte{1, 2, 3, 4},
	})
	signMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  []byte{1, 2, 3, 4},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         qbft.FirstRound,
		Input:         []byte{1, 2, 3, 4},
		PreparedRound: 2,
	})

	signMsgHeader, _ := signMsg.ToSignedMessageHeader()
	signMsgHeader2, _ := signMsg2.ToSignedMessageHeader()
	signMsgHeader3, _ := signMsg3.ToSignedMessageHeader()
	rcMsg.RoundChangeJustifications = []*qbft.SignedMessageHeader{
		signMsgHeader,
		signMsgHeader2,
		signMsgHeader3,
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
