package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangePrePreparedJustifications tests valid justified change round (prev prepared)
func RoundChangePrePreparedJustifications() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         10,
		PreparedRound: qbft.FirstRound,
	}, inputData)
	rcMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		signQBFTMsg,
		signQBFTMsg2,
		signQBFTMsg3,
	}
	rcMsgEncoded, _ := rcMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
	}

	return &tests.MsgSpecTest{
		Name:     "rc prev prepared justifications",
		Messages: msgs,
	}
}
