package decided

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// CurrentInstance tests a decided msg received for current running instance
func CurrentInstance() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide current instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
						MsgType:    qbft.ProposalMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
					}),

					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
						MsgType:    qbft.PrepareMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
					}),
					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
						MsgType:    qbft.PrepareMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
					}),
					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
						MsgType:    qbft.PrepareMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
					}),

					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
					}),
					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
					}),

					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.CommitMsgType,
							Height:     qbft.FirstHeight,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: []byte{1, 2, 3, 4},
				},
				ControllerPostRoot: "12fb900494537e891f774dacdfc1ea1ae9e11fdb099221eed157b0b517672cae",
			},
		},
	}
}
