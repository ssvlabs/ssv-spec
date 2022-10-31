package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateRoundChangePastInstance tests process round change msg for a previously decided instance
func LateRoundChangePastInstance() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	allMsgs := testingutils.DecidingMsgsForHeight(inputData, identifier, 5, ks)

	msgPerHeight := make(map[qbft.Height][]*types.Message)
	msgPerHeight[qbft.FirstHeight] = allMsgs[0:7]
	msgPerHeight[1] = allMsgs[7:14]
	msgPerHeight[2] = allMsgs[14:21]
	msgPerHeight[3] = allMsgs[21:28]
	msgPerHeight[4] = allMsgs[28:35]
	msgPerHeight[5] = allMsgs[35:42]

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		multiSignMsg := testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				Height: height,
				Round:  qbft.FirstRound,
				Input:  inputData,
			})
		return &tests.RunInstanceData{
			InputValue:         inputData,
			InputMessages:      msgPerHeight[height],
			SavedDecided:       multiSignMsg,
			BroadcastedDecided: multiSignMsg,
			DecidedVal:         inputData.Source,
			DecidedCnt:         1,
			ControllerPostRoot: postRoot,
		}
	}
	signMsgEncoded, _ := testingutils.SignQBFTMsg(ks.Shares[4], 4, &qbft.Message{
		Height: 4,
		Round:  qbft.FirstRound,
		Input:  inputData,
	}).Encode()

	return &tests.ControllerSpecTest{
		Name: "late proposal past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "d5d4696d29f1359a0f55292ba42dfd922993408529aa86926243df2221554c11"),
			instanceData(1, "457bc465febc4d1d626ad19d0f83621fbca5f0c2c6f9f3665292c602e615896a"),
			instanceData(2, "060c2a36313e8de4cfa530d7839945439e54344c3368eab749c61c5a76eb602c"),
			instanceData(3, "80f4ea4b56c6062724bc789eb3455c33650191e2c7f775f59b40b9fecc35f93b"),
			instanceData(4, "3003436a999f2fbd9d4c130591361243190fd3ab1da6d92463cbc832f8165abf"),
			instanceData(5, "e7a2324d9cbd69497455b50bde88cb47524b79b14653d024caf06ac7a2b28ba7"),
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
						Data: signMsgEncoded,
					},
				},
				ControllerPostRoot: "8368611d5315a0f2cfff2b8f02148f5badd52cda75f7cdb3878a0385faf99281",
			},
		},
	}
}
