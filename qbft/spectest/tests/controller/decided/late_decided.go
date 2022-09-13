package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateDecided tests processing a decided msg for a just decided instance
func LateDecided() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, testingutils.Testing4SharesSet())
	msgs = append(msgs, testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[4]},
		[]types.OperatorID{1, 2, 4},
		&qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: identifier[:],
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		}))
	return &tests.ControllerSpecTest{
		Name: "decide late decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				SavedDecided: testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
					}),
				BroadcastedDecided: testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
					}),
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         1,
				ControllerPostRoot: "aa402d7487719b17dde352e2ac602ba2c7d895e615ab12cd93d816f6c4fa0967",
			},
		},
	}
}
