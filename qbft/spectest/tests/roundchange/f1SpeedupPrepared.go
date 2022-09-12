package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1SpeedupPrepared tests catching up to higher rounds via f+1 speedup, other peers at the same round and already prepared
func F1SpeedupPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
		Input:  []byte{1, 2, 3, 4},
	})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
		Input:  []byte{1, 2, 3, 4},
	})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
		Input:  []byte{1, 2, 3, 4},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         10,
		Input:         []byte{1, 2, 3, 4},
		PreparedRound: 8,
	})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         10,
		Input:         []byte{1, 2, 3, 4},
		PreparedRound: 8,
	})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         10,
		Input:         []byte{1, 2, 3, 4},
		PreparedRound: 8,
	})
	proposalMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	proposalMsg10Round := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	})
	signQBFTMsg10RoundEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	signQBFTMsg10RoundEncoded2, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	signQBFTMsg10RoundEncoded3, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	signQBFTMsgFirstRoundEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	outputRcMsg10RoundEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  nil,
	}).Encode()

	prepareMsgHeader, _ := signQBFTMsg.ToSignedMessageHeader()
	prepareMsgHeader2, _ := signQBFTMsg2.ToSignedMessageHeader()
	prepareMsgHeader3, _ := signQBFTMsg3.ToSignedMessageHeader()

	prepareJustifications := []*qbft.SignedMessageHeader{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}
	rcMsg.RoundChangeJustifications = prepareJustifications
	rcMsg2.RoundChangeJustifications = prepareJustifications
	rcMsg3.RoundChangeJustifications = prepareJustifications

	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsgEncoded2, _ := rcMsg2.Encode()
	rcMsgEncoded3, _ := rcMsg3.Encode()

	rcMsgHeader, _ := rcMsg.ToSignedMessageHeader()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessageHeader()
	rcMsgHeader3, _ := rcMsg3.ToSignedMessageHeader()

	rcJustifications := []*qbft.SignedMessageHeader{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader3,
	}
	proposalMsg10Round.RoundChangeJustifications = rcJustifications
	proposalMsg10Round.ProposalJustifications = prepareJustifications
	proposalMsg10RoundEncoded, _ := proposalMsg10Round.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposalMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded3,
		},
		{ // this is the catchup round change sent by the node
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposalMsg10RoundEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signQBFTMsg10RoundEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signQBFTMsg10RoundEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signQBFTMsg10RoundEncoded3,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signQBFTMsg10RoundEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signQBFTMsg10RoundEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signQBFTMsg10RoundEncoded3,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "f+1 speed up (future prepared)",
		Pre:              pre,
		PostRoot:         "9e18563ecf2c03cc5181569df80b391c46499ba4b95392434ddb426ce64da3f8",
		InputMessagesSIP: msgs,
		OutputMessagesSIP: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signQBFTMsgFirstRoundEncoded,
			},
			{ // this is the catchup round change sent by the node
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
				Data: outputRcMsg10RoundEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposalMsg10RoundEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signQBFTMsg10RoundEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
				Data: signQBFTMsg10RoundEncoded,
			},
		},
	}
}
