package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// UnknownSigner tests a decided msg with an unknown signer
func UnknownSigner() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	multiSignMsg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 5},
		&qbft.Message{
			Height: 10,
			Round:  qbft.FirstRound,
			Input:  inputData,
		})
	multiSignMsgEncoded, _ := multiSignMsg.Encode()
	return &tests.ControllerSpecTest{
		Name: "decide unknown signer",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
						Data: multiSignMsgEncoded,
					},
				},
				ControllerPostRoot: "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
			},
		},
		ExpectedError: "invalid decided msg: decided msg signature invalid: unknown signer",
	}
}
