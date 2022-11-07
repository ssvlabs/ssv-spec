package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests a signed round change msg with an unknown signer
func UnknownSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	rcMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(5), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change unknown signer",
		Pre:            pre,
		PostRoot:       "a8b80879ebf2ecee42fddc69b67dd5f6adfd6aa8b7114246aec80ce1bfef513a",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "round change msg invalid: round change msg signature invalid: unknown signer",
	}
}
