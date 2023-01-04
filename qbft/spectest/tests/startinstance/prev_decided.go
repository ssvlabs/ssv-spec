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
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "start instance post future decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
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
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal:               []byte{1, 2, 3, 4},
					DecidedCnt:               1,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{qbft.FirstHeight, 10},
				},
				ControllerPostRoot: "6ec856c53c4febbeeb0d816b81a04425f5a7bdf107c7cf3d28a519c3fee6ce6e",
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal:               []byte{1, 2, 3, 4},
					DecidedCnt:               0,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{qbft.FirstHeight, 10},
				},
				ControllerPostRoot: "9255cf35e771cc91f931f1ed291137c7de4229d3fd2b46a638f1e214bd2a9a04",
			},
		},
	}
}
