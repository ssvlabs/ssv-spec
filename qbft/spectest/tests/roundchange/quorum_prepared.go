package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// QuorumPrepared tests a round change msg for prepared state
func QuorumPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: qbft.FirstRound,
	})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: qbft.FirstRound,
	})
	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})

	prepareMsgHeader, _ := signQBFTMsg.ToSignedMessage()
	prepareMsgHeader2, _ := signQBFTMsg2.ToSignedMessage()
	prepareMsgHeader3, _ := signQBFTMsg3.ToSignedMessage()

	prepareJustifications := []*qbft.SignedMessage{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}
	rcMsg.RoundChangeJustifications = prepareJustifications
	rcMsg3.RoundChangeJustifications = prepareJustifications

	rcMsgHeader, _ := rcMsg.ToSignedMessage()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessage()
	rcMsgHeader3, _ := rcMsg3.ToSignedMessage()

	proposalMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader3,
	}
	proposalMsg.ProposalJustifications = prepareJustifications

	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsgEncoded2, _ := rcMsg2.Encode()
	rcMsgEncoded3, _ := rcMsg3.Encode()
	proposalMsgEncoded, _ := proposalMsg.Encode()

	msgs := []*types.Message{
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

	return &tests.MsgProcessingSpecTest{
		Name:          "round change prepared",
		Pre:           pre,
		PostRoot:      "693f301963e027b305656d88af9eeb312f70216c49b16661a8ffce3fc6409e70",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposalMsgEncoded,
			},
		},
	}
}
