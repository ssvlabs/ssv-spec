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

	proposalMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
	}, pre.StartValue)
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
		Round:  9,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()
	proposalMsg.RoundChangeJustifications = []*qbft.SignedMessage{
		rcMsg.ToJustification(),
		rcMsg2.ToJustification(),
		rcMsg3.ToJustification(),
	}
	pre.State.ProposalAcceptedForCurrentRound = proposalMsg

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: prepareMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "prepare prev round",
		Pre:           pre,
		PostRoot:      "7c11857af9ccc2a3b05c2d379ff79976cae3f807af41390e560e14006c01b7da",
		InputMessages: msgs,
		ExpectedError: "invalid prepare msg: msg round wrong",
	}
}
