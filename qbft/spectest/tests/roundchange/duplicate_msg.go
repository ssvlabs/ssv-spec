package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate rc msg (first one inserted, second not)
func DuplicateMsg() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	prepareMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	prepareMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	prepareMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	changeRoundMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
	}, &qbft.Data{})
	changeRoundMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         5,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)

	changeRoundMsg2.RoundChangeJustifications = []*qbft.SignedMessage{
		prepareMsg,
		prepareMsg2,
		prepareMsg3,
	}

	changeRoundMsgEncoded, _ := changeRoundMsg.Encode()
	changeRoundMsgEncoded2, _ := changeRoundMsg2.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded2,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change duplicate msg",
		Pre:           pre,
		PostRoot:      "b8a30302b0e9a23b0b2a6b0d044bc7c0660efc7e85940e42973dfd084ba748b3",
		InputMessages: msgs,
	}
}
