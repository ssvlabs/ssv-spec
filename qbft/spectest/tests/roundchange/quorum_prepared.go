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
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         pre.StartValue,
		PreparedRound: qbft.FirstRound,
	})
	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: pre.StartValue.Root},
	})

	prepareJustifications := []*qbft.SignedMessage{
		signQBFTMsg, signQBFTMsg2, signQBFTMsg3,
	}
	rcMsg.RoundChangeJustifications = prepareJustifications
	rcMsg3.RoundChangeJustifications = prepareJustifications
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
		Name:          "round change prepared",
		Pre:           pre,
		PostRoot:      "fc31834e2b6b5ffe23ab33d723a7b060d4bca4f4fd2b946fd40e108b44682b67",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposalMsgEncoded,
			},
		},
	}
}
