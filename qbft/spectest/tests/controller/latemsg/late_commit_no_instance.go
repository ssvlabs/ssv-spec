package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateCommitNoInstance tests process commit msg for a previously decided instance (which is no longer part of stored instances)
func LateCommitNoInstance() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		multiSignMsg := testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				Height: height,
				Round:  qbft.FirstRound,
			}, inputData)
		multiSignMsgEncoded, _ := multiSignMsg.Encode()
		return &tests.RunInstanceData{
			InputValue: inputData,
			InputMessages: []*types.Message{
				{
					ID:   types.PopulateMsgType(identifier, types.DecidedMsgType),
					Data: multiSignMsgEncoded,
				},
			},
			SavedDecided:       multiSignMsg,
			DecidedVal:         inputData.Source,
			DecidedCnt:         1,
			ControllerPostRoot: postRoot,
		}
	}
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], 4, &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: inputData.Root}).Encode()

	return &tests.ControllerSpecTest{
		Name: "late commit no instance",
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
					{
						ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
						Data: signMsgEncoded,
					},
				},
				ControllerPostRoot: "4a184ea5a625621f73640e4f66209d38016e34030608c3afb0f6c2a33a38f886",
			},
		},
		ExpectedError: "instance not found",
	}
}
