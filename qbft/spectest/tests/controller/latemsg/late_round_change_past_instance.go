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
			ExpectedDecidedState: tests.DecidedState{
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
				DecidedVal: []byte{1, 2, 3, 4},
				DecidedCnt: 1,
			},

			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "late round change past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "0370be5066cbbf1efead61d9b182309afd989b3b720163f7029cbad79537eb4b"),
			instanceData(1, "72f925307e1d7af4d664a9e0339462a5c9d04bc17ceebe6084a506da89596d30"),
			instanceData(2, "411af5dd976dd829af4429d16c314e579a5ce98807122a4330071699e6cfdbb1"),
			instanceData(3, "21ef7d1175c6cc9c05ea489fc6ffa51215d0cce49a9bb4e9e14718e91803faf6"),
			instanceData(4, "1c2eebe8f3cc38f0063a3633a23827e0721298a4c4dfcc465bdc7f861f4bd75a"),
			instanceData(5, "2a95f62f8de66bc821c07c1476452aadacd45435b2a1a73c19f24d87460091ff"),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[4]},
						[]types.OperatorID{4},
						&qbft.Message{
							MsgType:    qbft.RoundChangeMsgType,
							Height:     4,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				ControllerPostRoot: "ff9a11fe4e4e17359000d5e2b063c05cca5b93d1f341e4c24a3bb16c4ace876b",
			},
		},
	}
}
