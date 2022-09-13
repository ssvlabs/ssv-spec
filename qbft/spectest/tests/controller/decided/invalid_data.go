package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// InvalidData tests decided data for which commitData.validate() != nil
func InvalidData() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide invalid data",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     10,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes(nil),
						}),
				},
				ControllerPostRoot: "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
			},
		},
		ExpectedError: "invalid decided msg: invalid decided msg: msgCommitData invalid: CommitData data is invalid",
	}
}
