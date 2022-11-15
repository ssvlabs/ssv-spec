package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureRoundPrevNotPrepared tests a proposal for future round, currently not prepared
func FutureRoundPrevNotPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = qbft.FirstRound

	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{})
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, pre.StartValue)
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()

	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg.ToJustification(),
		rcMsg2.ToJustification(),
		rcMsg3.ToJustification(),
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()

	return &tests.MsgProcessingSpecTest{
		Name:     "proposal future round prev not prepared",
		Pre:      pre,
		PostRoot: "41afd498ca82c22563a7a9726a22548ea7d991b25c69bbfbcedaec92022169c5",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposeMsgEncoded,
			},
		},
		OutputMessages: []*types.Message{
			{
				ID: types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded,
			},
		},
	}
}
