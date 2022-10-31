package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigner tests a proposal msg with > 1 signers
func MultiSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	signMsgInvalidEncoded, _ := testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2]}, []types.OperatorID{1, 2}, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  pre.StartValue,
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: signMsgInvalidEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "proposal multi signer",
		Pre:            pre,
		PostRoot:       "56cee2fd474513bc56851dfbb027366f6fc3f90fe8fec4081e993b69f84e2228",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "proposal invalid: proposal msg allows 1 signer",
	}
}
