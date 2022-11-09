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

	inputData := &qbft.Data{
		Root:   testingutils.TestAttesterConsensusDataRoot,
		Source: testingutils.TestAttesterConsensusDataByts,
	}
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
			}, inputData)
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
	}, inputData).Encode()

	return &tests.ControllerSpecTest{
		Name: "late change round past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "e7823a17225ee7f1163e71b0fc0b67df888cfe287f5ec7a6454ab105a402a998"),
			instanceData(1, "d3d2790b0746a32868d7ff7df0d115f84a0a33a81dde8be94efb482d0aa3c8a5"),
			instanceData(2, "dabcf05fba021ec68a0cb9f53b01b7a359c61bb24e78527eafe229e5c675b5b0"),
			instanceData(3, "542377f4dfe2349420d28ec87c653793a8844c5f63647b9388de34028f0c9248"),
			instanceData(4, "ea6870244b29a8d47a03cf859983582cf6b2855bfddc01b6801cfe104d576610"),
			instanceData(5, "d97b5ceafe67f4ebc2d708f0ba14ccc4220be4ae4802a75bd3d99e8821c3609f"),
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
						Data: signMsgEncoded,
					},
				},
				ControllerPostRoot: "3293219c0dc9b5b96459c4136ddc5129e028286383754928dbc1cd320e02fc36",
			},
		},
	}
}
