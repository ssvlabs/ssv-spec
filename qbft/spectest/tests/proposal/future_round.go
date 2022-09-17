package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FutureRound tests a proposal for state.ProposalAcceptedForCurrentRound != nil && signedProposal.Message.Round > state.Round
func FutureRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	pre.State.Round = qbft.FirstRound

	prepareMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
		Input:  []byte{1, 2, 3, 4},
	})
	prepareMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
		Input:  []byte{1, 2, 3, 4},
	})
	prepareMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
		Input:  []byte{1, 2, 3, 4},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	})
	prepareMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  []byte{1, 2, 3, 4},
	})

	prepareMsgHeader, _ := prepareMsg.ToSignedMessageHeader()
	prepareMsgHeader2, _ := prepareMsg2.ToSignedMessageHeader()
	prepareMsgHeader3, _ := prepareMsg3.ToSignedMessageHeader()

	rcMsgHeader, _ := rcMsg.ToSignedMessageHeader()
	rcMsgHeader2, _ := rcMsg2.ToSignedMessageHeader()
	rcMsgHeader3, _ := rcMsg3.ToSignedMessageHeader()

	proposeMsg.ProposalJustifications = []*qbft.SignedMessageHeader{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}

	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessageHeader{
		rcMsgHeader,
		rcMsgHeader2,
		rcMsgHeader3,
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: proposeMsgEncoded,
		},
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal future round prev prepared",
		Pre:           pre,
		PostRoot:      "9c725be6a189887dbe4f4917878ff579ea1d6842b02266015c7ef5b4e0703005",
		InputMessages: msgs,
		OutputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: prepareMsgEncoded,
			},
		},
	}
}
