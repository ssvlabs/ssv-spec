package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigner tests a round change msg with multiple signers
func MultiSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	rcMsgEncoded, _ := testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2]}, []types.OperatorID{types.OperatorID(1), types.OperatorID(2)}, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  2,
	}, &qbft.Data{Root: pre.StartValue.Root}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change multi signers",
		Pre:            pre,
		PostRoot:       "a8b80879ebf2ecee42fddc69b67dd5f6adfd6aa8b7114246aec80ce1bfef513a",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "round change msg invalid: round change msg allows 1 signer",
	}
}
