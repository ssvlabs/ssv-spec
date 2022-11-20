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
			DecidedVal:         []byte{1, 2, 3, 4},
			DecidedCnt:         1,
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "multi decide instances",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "2db1b6b59f13cd9b30f1afe09bcd62539c7061485435d8f134b86317d820e71d"),
			instanceData(1, "d2f56d69e871011f360c5e2733666f9389e02a7c30e81acacf258afab9d992cb"),
			instanceData(2, "3fb39692b356a635ed0be8a59c54ec64ab4dad4a45ace96f6c1cb8d3365f2d6b"),
			instanceData(3, "10a75c884203e4ba860bd9db6e55ba5989fdcfd45ce52a92d6b3e9a6adffb1dc"),
			instanceData(4, "2266f4d33838f251c22dcf787551bb6dd7381b689353b8147853338917dddf37"),
			instanceData(5, "b4842a41180885b9175a2e6b039e47b013864fe19113b6412991f41d917c3099"),
			instanceData(8, "c653ee5763862980dca63b99042191ad86e4aa1efe951f83513a9f5c1e8bb55b"),
			instanceData(9, "259a6315456362a70df4891e9ffd928dcec7fd30134effd2556aa1662121c4a7"),
			instanceData(10, "595ef559f1d141b2ac15fc2661c543c0db0cfe39fb18d38f9a35dca0e5d07055"),
		},
	}
}
