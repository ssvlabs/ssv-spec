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
	}, pre.StartValue)
	pre.State.Round = qbft.FirstRound

	prepareMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
	}, &qbft.Data{Root: pre.StartValue.Root})
	prepareMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
	}, &qbft.Data{Root: pre.StartValue.Root})
	prepareMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  8,
	}, &qbft.Data{Root: pre.StartValue.Root})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{})
	rcMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{})
	rcMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{})
	prepareMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	proposeMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, pre.StartValue)

	proposeMsg.ProposalJustifications = []*qbft.SignedMessage{
		prepareMsg,
		prepareMsg2,
		prepareMsg3,
	}

	proposeMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg.ToJustification(),
		rcMsg2.ToJustification(),
		rcMsg3.ToJustification(),
	}
	proposeMsgEncoded, _ := proposeMsg.Encode()
	return &tests.MsgProcessingSpecTest{
		Name:     "proposal future round prev prepared",
		Pre:      pre,
		PostRoot: "5814e350322c5cc06b3dd99b77c6b97608de17063d67fab0c9fce26aa40a59d8",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposeMsgEncoded,
			},
		},
		OutputMessages: []*types.Message{
			{
				ID: types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
				Data: prepareMsgEncoded,
			},
		},
	}
}
