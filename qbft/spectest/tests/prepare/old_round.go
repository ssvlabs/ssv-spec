package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// OldRound tests prepare for signedProposal.Message.Round < state.Round
func OldRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 10
	rcMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      10,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      10,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      10,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      10,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, rcMsgs, nil),
	})

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      9,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare prev round",
		Pre:           pre,
		PostRoot:      "3a85ee7abc9275684534414a639bd4d968f66bb319214da057a47ab8f9a881b4",
		InputMessages: msgs,
		ExpectedError: "invalid prepare msg: msg round wrong",
	}
}
