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
			instanceData(3, "29648f3c50aec14f76731c0650a45f64373da08010ec4cee7236d1156558a213"),
			instanceData(7, "85bd7b4f02b2d1b29eef0f171c3cc5e208a2f39bf0f00bc2a3bd955dab2d4389"),
			instanceData(11, "e7ff49c725608be250eebf277e5d88c6408750e527759db5546a7abc0a186c98"),
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
				ControllerPostRoot: "e0204ddedc95ffe4360ae469877cb3a28745fe8e74a681b2595243aed5489b6c",
			},
		},
		ExpectedError: "instance not found",
	}
}
