package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigners tests future msg with multiple signers
func MultiSigners() *ControllerSyncSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	multiSignMsgEncoded, _ := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: 2,
			Round:  qbft.FirstRound,
			Input:  &qbft.Data{Root: inputData.Root},
		}).Encode()

	return &ControllerSyncSpecTest{
		Name: "future msgs multiple signers",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusPrepareMsgType),
				Data: multiSignMsgEncoded,
			},
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
		ExpectedError:        "invalid future msg: allows 1 signer",
	}
}
