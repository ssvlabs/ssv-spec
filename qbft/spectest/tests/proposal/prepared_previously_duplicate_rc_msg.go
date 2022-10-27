package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreparedPreviouslyDuplicateRCMsg tests a proposal for > 1 round, prepared previously with quorum of round change but 2 are duplicates (shouldn't find quorum)
func PreparedPreviouslyDuplicateRCMsg() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	prepareMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	prepareMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	prepareMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
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
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})

	prepareMsgHeader, _ := prepareMsg.ToSignedMessage()
	prepareMsgHeader2, _ := prepareMsg2.ToSignedMessage()
	prepareMsgHeader3, _ := prepareMsg3.ToSignedMessage()

	justifications := []*qbft.SignedMessage{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}

	rcMsg.RoundChangeJustifications = justifications
	rcMsg2.RoundChangeJustifications = justifications

	rcMsgHeader, _ := rcMsg.ToSignedMessage()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessage()

	proposeMsg.ProposalJustifications = justifications
	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader2,
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "duplicate rc msg justification (prepared)",
		Pre:            pre,
		PostRoot:       "3e721f04a2a64737ec96192d59e90dfdc93f166ec9a21b88cc33ee0c43f2b26a",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "proposal invalid: proposal not justified: change round has not quorum",
	}
}
