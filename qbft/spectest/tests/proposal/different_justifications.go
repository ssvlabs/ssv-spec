package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DifferentJustifications tests a proposal for > 1 round, prepared previously with rc justification prepares at different heights (tests the highest prepared calculation)
func DifferentJustifications() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 3

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
	signQBFTMsg4 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	signQBFTMsg5 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	signQBFTMsg6 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         3,
		Input:         &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: qbft.FirstRound,
	})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         3,
		Input:         &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: 2,
	})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  3,
	})
	prepareMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  3,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	}).Encode()
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  3,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})

	prepareMsgHeader, _ := signQBFTMsg.ToSignedMessage()
	prepareMsgHeader2, _ := signQBFTMsg2.ToSignedMessage()
	prepareMsgHeader3, _ := signQBFTMsg3.ToSignedMessage()

	prepareMsgHeader4, _ := signQBFTMsg4.ToSignedMessage()
	prepareMsgHeader5, _ := signQBFTMsg5.ToSignedMessage()
	prepareMsgHeader6, _ := signQBFTMsg6.ToSignedMessage()

	prepareJustifications := []*qbft.SignedMessage{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}
	prepareJustifications2 := []*qbft.SignedMessage{
		prepareMsgHeader4,
		prepareMsgHeader5,
		prepareMsgHeader6,
	}
	rcMsg.RoundChangeJustifications = prepareJustifications
	rcMsg2.RoundChangeJustifications = prepareJustifications2

	rcMsgHeader, _ := rcMsg.ToSignedMessage()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessage()
	rcMsgHeader3, _ := rcMsg3.ToSignedMessage()

	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader3,
	}
	proposeMsg.ProposalJustifications = prepareJustifications2
	proposeMsgEncoded, _ := proposeMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "different proposal round change justification",
		Pre:           pre,
		PostRoot:      "bfce5f9bd15d7dfb986713af965ea5b090909f688d8a142f2f03cb0edcc7b853",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: prepareMsgEncoded,
			},
		},
	}
}
