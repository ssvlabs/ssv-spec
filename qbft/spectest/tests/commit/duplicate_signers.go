package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// DuplicateSigners tests a multi signer commit msg with duplicate signers
func DuplicateSigners() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, pre.StartValue)
	commit := testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2]}, []types.OperatorID{1, 2}, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: pre.StartValue.Root})
	commit.Signers = []types.OperatorID{1, 1}
	commitEncoded, _ := commit.Encode()

	return &tests.MsgProcessingSpecTest{
		Name:     "duplicate signers",
		Pre:      pre,
		PostRoot: "b19c1eb98711766bb8b1e2857cafc8b83e9584c7fd8b9a3a81fb6947df1497f6",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
				Data: commitEncoded,
			},
		},
		OutputMessages: []*types.Message{},
		ExpectedError:  "invalid signed message: non unique signer",
	}
}
