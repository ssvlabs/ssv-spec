package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1Speedup tests catching up to higher rounds via f+1 speedup, other peers are all at the same round
func F1Speedup() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	})
	changeRoundMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  nil,
	})
	changeRoundMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  nil,
	})
	changeRoundMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  nil,
	})
	signQBFTMsgFirstRound := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	proposalMsg10 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	})

	changeRoundMsgHeader, _ := changeRoundMsg.ToSignedMessageHeader()
	changeRoundMsgHeader2, _ := changeRoundMsg2.ToSignedMessageHeader()
	changeRoundMsgHeader3, _ := changeRoundMsg3.ToSignedMessageHeader()

	proposalMsg10.RoundChangeJustifications = []*qbft.SignedMessageHeader{
		changeRoundMsgHeader,
		changeRoundMsgHeader2,
		changeRoundMsgHeader3,
	}

	changeRoundMsgEncoded, _ := changeRoundMsg.Encode()
	changeRoundMsgEncoded2, _ := changeRoundMsg2.Encode()
	changeRoundMsgEncoded3, _ := changeRoundMsg3.Encode()
	signQBFTMsgEncoded, _ := signQBFTMsg.Encode()
	signQBFTMsgEncoded2, _ := signQBFTMsg2.Encode()
	signQBFTMsgEncoded3, _ := signQBFTMsg3.Encode()
	signQBFTMsgFirstRoundEncoded, _ := signQBFTMsgFirstRound.Encode()
	proposalMsg10Encoded, _ := proposalMsg10.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: signQBFTMsgFirstRoundEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded3,
		},
		{ // this is the catchup round change sent by the node
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposalMsg10Encoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signQBFTMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signQBFTMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signQBFTMsgEncoded3,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signQBFTMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signQBFTMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signQBFTMsgEncoded3,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "f+1 speed up",
		Pre:              pre,
		PostRoot:         "ce506da81bac03ee9118b57e36fc350748ab280a27062fba882c4c9ba603d07c",
		InputMessagesSIP: msgs,
		OutputMessagesSIP: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signQBFTMsgFirstRoundEncoded,
			},
			{ // this is the catchup round change sent by the node
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
				Data: changeRoundMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposalMsg10Encoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signQBFTMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
				Data: signQBFTMsgEncoded,
			},
		},
	}
}
