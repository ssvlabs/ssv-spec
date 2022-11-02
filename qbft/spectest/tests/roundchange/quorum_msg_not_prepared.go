package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// QuorumMsgNotPrepared tests LIVENESS where the rc msg making a quorum for round change is not prev prepared (other are) which can cause a liveness issue with isReceivedProposalJustification
func QuorumMsgNotPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         pre.StartValue,
		PreparedRound: qbft.FirstRound,
	})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{},
	})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{},
	})
	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  pre.StartValue,
	})

	prepareJustifications := []*qbft.SignedMessage{
		signQBFTMsg, signQBFTMsg2, signQBFTMsg3,
	}
	rcMsg.RoundChangeJustifications = prepareJustifications
	proposalMsg.ProposalJustifications = prepareJustifications
	proposalMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg, rcMsg2, rcMsg3,
	}

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
		Name:          "round change quorum msg not prepared",
		Pre:           pre,
		PostRoot:      "6cae126d16d422505088c0affe82d726fb1a512c6693c0d5b5305a52136878a5",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposalMsgEncoded,
			},
		},
	}
}
