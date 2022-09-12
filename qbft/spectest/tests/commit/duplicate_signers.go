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
	// TODO<olegshmuelov> use instance identifier
	baseMsgId := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	commit := testingutils.MultiSignQBFTMsg([]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2]}, []types.OperatorID{1, 2}, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	commit.Signers = []types.OperatorID{1, 1}
	commitEncoded, _ := commit.Encode()

	return &tests.MsgProcessingSpecTest{
		Name:     "duplicate signers",
		Pre:      pre,
		PostRoot: "20a518595c0dbe81ccc7f340f142e77ecfba0e0a93fe0d10325fe607f2e0b1eb",
		InputMessagesSIP: []*types.Message{
			{
				ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
				Data: commitEncoded,
			},
		},
		OutputMessagesSIP: []*types.Message{},
		ExpectedError:     "invalid signed message: non unique signer",
	}
}
