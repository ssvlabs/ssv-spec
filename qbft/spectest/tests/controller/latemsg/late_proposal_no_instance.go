package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateProposalNoInstance tests process proposal msg for a previously decided instance (which is no longer part of stored instances)
func LateProposalNoInstance() *tests.ControllerSpecTest {
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
			ExpectedDecidedState: tests.DecidedState{
				DecidedVal:               []byte{1, 2, 3, 4},
				DecidedCnt:               1,
				CalledSyncDecidedByRange: true,
				DecidedByRangeValues:     [2]qbft.Height{height - 3, height},
			},
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "late proposal no instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(3, "ef701ee47bd9cf0c877ebe7148bc210803da594e12837c632fd9c56f658956ff"),
			instanceData(7, "d37a7c38cdea9f71111410ca3d6db6b9ca5746caa8833f722d2cf3c552ea735b"),
			instanceData(11, "78bce27afd40e40c155d5701823a983a3d20ba6bae6f74dafd228f9b4f434082"),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1]},
						[]types.OperatorID{1},
						&qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     2,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				ExpectedDecidedState: tests.DecidedState{
					CalledSyncDecidedByRange: true, // leftovers from the previous calls
					DecidedByRangeValues:     [2]qbft.Height{8, 11},
				},
				ControllerPostRoot: "21a7699ab552e583595d3138b87c99aa82cae149e592c7607b964da4b793f96a",
			},
		},
		ExpectedError: "instance not found",
	}
}
