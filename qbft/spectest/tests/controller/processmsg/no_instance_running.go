package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// NoInstanceRunning tests a process msg for height in which there is no running instance
func NoInstanceRunning() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "no instance running",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     50,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
					testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
						MsgType:    qbft.ProposalMsgType,
						Height:     2,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
					}),
				},

				ExpectedDecidedState: tests.DecidedState{
					DecidedVal:               []byte{1, 2, 3, 4},
					DecidedCnt:               1,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{0, 50},
				},
				ControllerPostRoot: "e80301d8150da48514ea573693b0f5aaf07e4e26727c66068ddfa208ad58722d",
			},
		},
		ExpectedError: "instance not found",
	}
}
