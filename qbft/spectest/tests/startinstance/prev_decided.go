package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PostFutureDecided tests starting a new instance after deciding with future decided msg
func PostFutureDecided() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "start instance post future decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     10,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				SavedDecided: testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     10,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
					}),
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "cf79aafb1c5b30c9c737026d3263934ee31da72db9b68bedc21c5e0f7def69b0",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				DecidedCnt:         0,
				ControllerPostRoot: "e7091248ed58bffb5751b0006a1f9c2e79760268e6f4d2c4efe6c30c792dc461",
			},
		},
	}
}
