package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreparedPreviouslyJustification tests a proposal for > 1 round, prepared previously with quorum of round change msgs justification
func PreparedPreviouslyJustification() *tests.MsgProcessingSpecTest {
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
		Height:        qbft.FirstHeight,
		Round:         2,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		PreparedRound: qbft.FirstRound,
	}, pre.StartValue)
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue).Encode()
	proposeMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, pre.StartValue)
	prepareRound2MsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()

	justifications := []*qbft.SignedMessage{
		signQBFTMsg,
		signQBFTMsg2,
		signQBFTMsg3,
	}

	rcMsg.RoundChangeJustifications = justifications
	rcMsg2.RoundChangeJustifications = justifications
	rcMsg3.RoundChangeJustifications = justifications

	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsgEncoded2, _ := rcMsg2.Encode()
	rcMsgEncoded3, _ := rcMsg3.Encode()

	prepareMsgEncoded, _ := signQBFTMsg.Encode()
	prepareMsgEncoded2, _ := signQBFTMsg2.Encode()
	prepareMsgEncoded3, _ := signQBFTMsg3.Encode()

	proposeMsg2.ProposalJustifications = justifications
	proposeMsg2.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg.ToJustification(),
		rcMsg2.ToJustification(),
		rcMsg3.ToJustification(),
	}
	proposeMsgEncoded2, _ := proposeMsg2.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: prepareMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: prepareMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: prepareMsgEncoded3,
		},
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
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded2,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "previously prepared proposal",
		Pre:           pre,
		PostRoot:      "2949f309dc07bc54627ec1cec2dd4a0fbf870729af9f6ce7cf10831b24888355",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: prepareMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
				Data: prepareMsgEncoded,
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
				ID: types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: prepareRound2MsgEncoded,
			},
		},
	}
}
