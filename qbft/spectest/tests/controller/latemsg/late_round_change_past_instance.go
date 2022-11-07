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
			instanceData(qbft.FirstHeight, "df9b2787df60e1e15b0c840410592d27803d44cd5fb086cfa8fc23181cea6293"),
			instanceData(1, "34d48526093a91455d5897e7787a80ffa46eec1f0d6bf54792a6163abf08eb0d"),
			instanceData(2, "c048c170bef823691ddfd12bbaea54339b26a12c0e9b92d25434f10a564c9467"),
			instanceData(3, "aa247f887c91f4442d83d545d383559c9c026e67e24509cab02fe60ec56470f9"),
			instanceData(4, "f1cdaff9b73dab929b270fcc592f71bb5f9f12c770ca89d95883f905a1388dd1"),
			instanceData(5, "481a19848af8107171f1b02b22c6363c54962c2b311d1e89c882a002b7296a02"),
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
						Data: signMsgEncoded,
					},
				},
				ControllerPostRoot: "555c8f4c00a06f5fb10940bb3f3c597841c32931a8a2d56d46a5951fc13fd4ec",
			},
		},
	}
}
