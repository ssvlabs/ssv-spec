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
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: qbft.FirstRound,
	})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: qbft.FirstRound,
	})
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	proposalFirstRound, _ := signQBFTMsg.Encode()
	outputPrepareMsgEncoded, _ := proposeMsg.Encode()

	prepareMsgEncoded, _ := signQBFTMsg.Encode()
	prepareMsgEncoded2, _ := signQBFTMsg2.Encode()
	prepareMsgEncoded3, _ := signQBFTMsg3.Encode()

	prepareMsgHeader, _ := signQBFTMsg.ToSignedMessage()
	prepareMsgHeader2, _ := signQBFTMsg2.ToSignedMessage()
	prepareMsgHeader3, _ := signQBFTMsg3.ToSignedMessage()

	justifications := []*qbft.SignedMessage{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}

	rcMsg.RoundChangeJustifications = justifications
	rcMsg2.RoundChangeJustifications = justifications
	rcMsg3.RoundChangeJustifications = justifications

	rcMsgEncoded, _ := rcMsg.Encode()
	rcMsgEncoded2, _ := rcMsg2.Encode()
	rcMsgEncoded3, _ := rcMsg3.Encode()

	rcMsgHeader, _ := rcMsg.ToSignedMessage()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessage()
	rcMsgHeader3, _ := rcMsg3.ToSignedMessage()

	proposeMsg.ProposalJustifications = justifications
	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader3,
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposalFirstRound,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: prepareMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: prepareMsgEncoded2,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
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
			Data: proposeMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "previously prepared proposal",
		Pre:           pre,
		PostRoot:      "b67f687963630a773f1dd31848b4b78938a51147fc5d7c3b7b0bb61d8008c5b7",
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
				Data: prepareMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposeMsgEncoded,
			},
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: outputPrepareMsgEncoded,
			},
		},
	}
}
