package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// FullDecided tests process msg and first time deciding
func FullDecided() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "first decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: testingutils.DecidingMsgsForHeight(&qbft.Data{
					Root:   [32]byte{1, 2, 3, 4},
					Source: []byte{1, 2, 3, 4},
				}, identifier, qbft.FirstHeight, ks),
				DecidedVal: []byte{1, 2, 3, 4},
				DecidedCnt: 1,
				SavedDecided: testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						//MsgType:    qbft.CommitMsgType,
						Height: qbft.FirstHeight,
						Round:  qbft.FirstRound,
						//Identifier: identifier[:],
						//Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						Input: &qbft.Data{
							Root:   [32]byte{1, 2, 3, 4},
							Source: nil,
						},
					}),
				BroadcastedDecided: testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						//MsgType:    qbft.CommitMsgType,
						Height: qbft.FirstHeight,
						Round:  qbft.FirstRound,
						//Identifier: identifier[:],
						//Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						Input: &qbft.Data{
							Root:   [32]byte{1, 2, 3, 4},
							Source: nil,
						},
					}),
				ControllerPostRoot: "aa402d7487719b17dde352e2ac602ba2c7d895e615ab12cd93d816f6c4fa0967",
			},
		},
	}
}
