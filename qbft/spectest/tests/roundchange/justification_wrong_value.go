package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// JustificationWrongValue tests a single prepare justification with value != prepared value
func JustificationWrongValue() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		PreparedRound: qbft.FirstRound,
	}, &qbft.Data{Root: [32]byte{1, 2, 3, 3}, Source: []byte{1, 2, 3, 3}})

	rcMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		signQBFTMsg, signQBFTMsg2, signQBFTMsg3,
	}
	rcMsgEncoded, _ := rcMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change justification wrong value",
		Pre:            pre,
		PostRoot:       "a8b80879ebf2ecee42fddc69b67dd5f6adfd6aa8b7114246aec80ce1bfef513a",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "round change msg invalid: round change justification invalid: prepare data != proposed data",
	}
}
