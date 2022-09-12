package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsgQuorum tests a duplicate rc msg that is added to container but shouldn't result in a quorum
func DuplicateMsgQuorum() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	prepareMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	prepareMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	prepareMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	changeRoundMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  nil,
	})
	changeRoundMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         []byte{1, 2, 3, 4},
		PreparedRound: qbft.FirstRound,
	})
	changeRoundMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  nil,
	})

	prepareMsgHeader, _ := prepareMsg.ToSignedMessageHeader()
	prepareMsgHeader2, _ := prepareMsg2.ToSignedMessageHeader()
	prepareMsgHeader3, _ := prepareMsg3.ToSignedMessageHeader()

	prepareJustifications := []*qbft.SignedMessageHeader{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}
	changeRoundMsg2.RoundChangeJustifications = prepareJustifications

	changeRoundMsgEncoded, _ := changeRoundMsg.Encode()
	changeRoundMsgEncoded2, _ := changeRoundMsg2.Encode()
	changeRoundMsgEncoded3, _ := changeRoundMsg3.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: changeRoundMsgEncoded3,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "round change duplicate msg",
		Pre:              pre,
		PostRoot:         "14e50f11f388e9ecf1ad7c5a18bd764fc30613b35a00b27a2f8ffe813925cec8",
		InputMessagesSIP: msgs,
		OutputMessages:   []*qbft.SignedMessage{},
	}
}
