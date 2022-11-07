package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1DifferentFutureRounds tests f+1 speedup with one rc prev prepared
func F1DifferentFutureRounds() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
	}, &qbft.Data{})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         10,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)

	rcMsg2.RoundChangeJustifications = []*qbft.SignedMessage{
		signQBFTMsg, signQBFTMsg2, signQBFTMsg3,
	}

	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsgEncoded2, _ := rcMsg2.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded2,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change f+1 prepared",
		Pre:           pre,
		PostRoot:      "a118a0b79b839acf6d3c4b9bce674ee4e7cced87eaf93cd0f4ba505e9805adcd",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
				Data: rcMsgEncoded,
			},
		},
	}
}
