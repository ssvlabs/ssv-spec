package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoInstanceRunning tests a process msg for height in which there is no running instance
func NoInstanceRunning() *tests.ControllerSpecTest {
	//identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: 2,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	return &tests.ControllerSpecTest{
		Name: "no instance running",
		RunInstanceData: []struct {
			InputValue    []byte
			InputMessages []*types.Message
			Decided       bool
			DecidedVal    []byte
			DecidedCnt    uint
			SavedDecided  *qbft.SignedMessage
		}{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
						Data: signMsgEncoded,
					}},
				Decided:    false,
				DecidedVal: nil,
			},
		},
		ExpectedError: "instance not found",
	}
}
