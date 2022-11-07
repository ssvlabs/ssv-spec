package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSignature tests a single prepare received with a wrong signature
func WrongSignature() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue)

	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusPrepareMsgType),
			Data: signMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "prepare wrong sig",
		Pre:           pre,
		PostRoot:      "b19c1eb98711766bb8b1e2857cafc8b83e9584c7fd8b9a3a81fb6947df1497f6",
		InputMessages: msgs,
		ExpectedError: "invalid prepare msg: prepare msg signature invalid: failed to verify signature",
	}
}
