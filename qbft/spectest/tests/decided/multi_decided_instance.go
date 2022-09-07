package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiDecidedInstances tests deciding multiple instances
func MultiDecidedInstances() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue: []byte{1, 2, 3, 4},
			InputMessages: []*qbft.SignedMessage{
				testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     height,
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
					Height:     height,
					Round:      qbft.FirstRound,
					Identifier: identifier[:],
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				}),
			BroadcastedDecided: testingutils.MultiSignQBFTMsg(
				[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
				[]types.OperatorID{1, 2, 3},
				&qbft.Message{
					MsgType:    qbft.CommitMsgType,
					Height:     height,
					Round:      qbft.FirstRound,
					Identifier: identifier[:],
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				}),
			DecidedVal:         []byte{1, 2, 3, 4},
			DecidedCnt:         1,
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "multi decide instances",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "19e63c3c0d763de50f31c3d41dcf80af0c68f8ab659ff02239eb81f1ed757fef"),
			instanceData(1, "34c6fa58cd4bea9d0e0ce22b028cbf092ebc41b4961f83877a9a2bfa2ccc7dad"),
			instanceData(2, "bae57bcbf52062426a23f050f970c4f4c328f07208aba8a6200e5da2e0e77014"),
			instanceData(3, "733a85040e9a3820f9be77c06726bdaa28b0865023e32c44aac286e4cc0499cf"),
			instanceData(4, "5aa8f71a466fcb7a7761e155c3c2dc835ef5389b6e7f7feca7d53e3975fcc4aa"),
			instanceData(5, "7686425bc58541fc21b8309fb8ec64068bdda53cf41c8aea5dabec62d0544621"),
			instanceData(6, "6fbd7892e4e5c38c33b3fb81852f6cc120821312e48ff16da0e9e835608e6d78"),
			instanceData(7, "827822d44f5ee786b24dec5cf1d0ce155271cc0731e24b0c2d1c8e8609cd40f8"),
			instanceData(8, "12e8be02534cdf8653813527b5a7f6b976ad779375d6c854562636844b8a8f53"),
			instanceData(9, "aacc8eb183b778530a1b5d53f208ea79aaf46596265e0d4b2ea3f102c17095f3"),
			instanceData(10, "db4d7f1e60e109e207b3bd35fe7caa94a92253c9538d39572954baf70cf77550"),
		},
	}
}
