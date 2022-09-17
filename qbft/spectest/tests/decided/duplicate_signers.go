package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// DuplicateSigners tests a decided msg with duplicate signers
func DuplicateSigners() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	commit := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  []byte{1, 2, 3, 4},
		})
	commit.Signers = []types.OperatorID{1, 1, 3}
	commitMsgEncoded, _ := commit.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: commitMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "decided duplicate signers",
		Pre:            pre,
		PostRoot:       "be41977d818071451988105377df7c5ccf89ecc05ddf033b7b3b83d89f52d530",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "invalid signed message: non unique signer",
	}
}
