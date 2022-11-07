package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotPreparedPreviouslyJustification tests a proposal for > 1 round, not prepared previously
func NotPreparedPreviouslyJustification() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 5

	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
	}, &qbft.Data{})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
	}, &qbft.Data{})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
	}, &qbft.Data{})
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
	}, pre.StartValue)
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()

	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg,
		rcMsg2,
		rcMsg3,
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "proposal justification (not prepared)",
		Pre:           pre,
		PostRoot:      "36da1dbc06caceaaef99d529a8faecdfd96a26695b29111898ca724fda7817e9",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded,
			},
		},
	}
}
