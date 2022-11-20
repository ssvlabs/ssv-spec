package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LatePreparePastInstance tests process prepare msg for a previously decided instance
func LatePreparePastInstance() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	allMsgs := testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], 5, ks)
	msgPerHeight := make(map[qbft.Height][]*qbft.SignedMessage)
	msgPerHeight[qbft.FirstHeight] = allMsgs[0:7]
	msgPerHeight[1] = allMsgs[7:14]
	msgPerHeight[2] = allMsgs[14:21]
	msgPerHeight[3] = allMsgs[21:28]
	msgPerHeight[4] = allMsgs[28:35]
	msgPerHeight[5] = allMsgs[35:42]

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue:    []byte{1, 2, 3, 4},
			InputMessages: msgPerHeight[height],
			SavedDecided: testingutils.MultiSignQBFTMsg(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				&qbft.Message{
					MsgType:    qbft.CommitMsgType,
					Height:     height,
					Round:      qbft.FirstRound,
					Identifier: identifier[:],
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				}),
			BroadcastedDecided: testingutils.MultiSignQBFTMsg(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
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
		Name: "late prepare past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "f91546f051287e118a5b22ef4750062dc5d41fca0f5106cddbcd76447161ba88"),
			instanceData(1, "b8ca18cd97642b125293799d7c79a9b99419bd5779cf852b55baaf1d98cd0e35"),
			instanceData(2, "cc4f29db1f8055cfb64198c9951a660dd002bb22d7a585ef86d1906048701a80"),
			instanceData(3, "8059cec4734df8c9d6f3c38581c757093f14c820ff3c09f5b0f75e74dae69d1a"),
			instanceData(4, "3003436a999f2fbd9d4c130591361243190fd3ab1da6d92463cbc832f8165abf"),
			instanceData(5, "b246211aabb7a12cf76f160833ba02f53883654b8363a68ce4fe1994b4301034"),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[4]},
						[]types.OperatorID{4},
						&qbft.Message{
							MsgType:    qbft.PrepareMsgType,
							Height:     4,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				ControllerPostRoot: "5657c17279f7ad4f2f66e8331d7e373e47c43f04cf3a24c309e1b734627b2750",
			},
		},
	}
}
