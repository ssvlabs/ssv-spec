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
	signMsgEncodedFirstRound, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  nil,
	})
	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  nil,
	})
	rcMsgEncoded2, _ := rcMsg2.Encode()
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  nil,
	})
	rcMsgEncoded3, _ := rcMsg3.Encode()

	rcMsgHeader, _ := rcMsg.ToSignedMessageHeader()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessageHeader()
	rcMsgHeader3, _ := rcMsg3.ToSignedMessageHeader()

	signMsg2Round := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  []byte{1, 2, 3, 4},
	})
	signMsg2Round2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  []byte{1, 2, 3, 4},
	})
	signMsg2Round3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  []byte{1, 2, 3, 4},
	})
	signMsg2RoundEncoded, _ := signMsg2Round.Encode()
	signMsg2RoundEncoded2, _ := signMsg2Round2.Encode()
	signMsg2RoundEncoded3, _ := signMsg2Round3.Encode()
	signMsg2Round.RoundChangeJustifications = []*qbft.SignedMessageHeader{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader3,
	}
	signMsg2RoundEncodedWithJust, _ := signMsg2Round.Encode()

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
			Data: signMsgEncodedFirstRound,
		},
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs,
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: signMsg2RoundEncodedWithJust,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsg2RoundEncoded,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsg2RoundEncoded2,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsg2RoundEncoded3,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsg2RoundEncoded,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsg2RoundEncoded2,
		},
		&types.Message{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsg2RoundEncoded3,
		},
	)

	return &tests.MsgProcessingSpecTest{
		Name:             "round change happy flow",
		Pre:              pre,
		PostRoot:         "04ef76c9b07f2f02f8cad332bd2ed331985d214be3a461e8f996d6e771901a8f",
		InputMessagesSIP: msgs,
		OutputMessagesSIP: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signMsgEncodedFirstRound,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
				Data: rcMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: signMsg2RoundEncodedWithJust,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: signMsg2RoundEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
				Data: signMsg2RoundEncoded,
			},
		},
	}
}
