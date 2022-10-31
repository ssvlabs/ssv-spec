package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateRoundChangeNoInstance tests process round change msg for a previously decided instance (which is no longer part of stored instances)
func LateRoundChangeNoInstance() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue: inputData,
			InputMessages: []*types.Message{
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
			DecidedVal:         inputData.Source,
			DecidedCnt:         1,
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "late round change no instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "e719fcf29fb1acbd0dbf0843f61c2d037463e6fe7c4e21e78ad28825c4c56e41"),
			instanceData(1, "8472719bee2bb179963cc3aee61f054042de5f40f45cce557dbb8030ea5f32ca"),
			instanceData(2, "03fbaa2914cec6145e0799ea822489048ec14cc3c9492b41f13f152440cd7fc5"),
			instanceData(3, "e14a8e324b78dda8469df11a8d3551f6fac6163d6458d51fee7853c0f17d5835"),
			instanceData(4, "bd54d2ab1e0b949dd45a3e2b16f211d006371ae99a96f7a8d74fa98128447bf9"),
			instanceData(5, "3102e5048ac659bd0bbc8efc4eab07154f357a1225c8b1ef0681e054310f3bb3"),
			instanceData(8, "c880aac62523021de564af42efb73d9da75bafa087714a03291bd5e6a6ba3acd"),
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[4]},
						[]types.OperatorID{4},
						&qbft.Message{
							MsgType:    qbft.RoundChangeMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				ControllerPostRoot: "69b8be89bbd2b48e81644c689709d6ce0f43f898a2f7ef3b888d84a0cb264db7",
			},
		},
		ExpectedError: "instance not found",
	}
}
