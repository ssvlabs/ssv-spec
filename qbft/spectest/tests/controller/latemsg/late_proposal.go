package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateProposal tests process late proposal msg for an instance which just decided
func LateProposal() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	msgs := testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, ks)
	msgs = append(msgs, testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: identifier[:],
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	}))

	return &tests.ControllerSpecTest{
		Name: "late proposal",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: []byte{1, 2, 3, 4},
					DecidedCnt: 1,
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
				},
				ControllerPostRoot: "0370be5066cbbf1efead61d9b182309afd989b3b720163f7029cbad79537eb4b",
			},
		},
		ExpectedError: "could not process msg: proposal invalid: proposal is not valid with current state",
	}
}
