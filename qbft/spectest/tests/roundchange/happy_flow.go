package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue).Encode()
	proposeMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, pre.StartValue)
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{})

	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	signMsgEncoded4, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	proposeMsg2.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg, rcMsg2, rcMsg3,
	}
	proposeMsgEncoded2, _ := proposeMsg2.Encode()
	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsgEncoded2, _ := rcMsg2.Encode()
	rcMsgEncoded3, _ := rcMsg3.Encode()

	rcMsgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded3,
		},
	}

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs,
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded2,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded2,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded3,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded2,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded3,
		},
	)

	return &tests.MsgProcessingSpecTest{
		Name:          "round change happy flow",
		Pre:           pre,
		PostRoot:      "ea7943b99324f66e94e82a287aaa449f72850671b1fef51cb480bb81a62ec1d8",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded4,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
				Data: rcMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposeMsgEncoded2,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
				Data: signMsgEncoded,
			},
		},
	}
}
