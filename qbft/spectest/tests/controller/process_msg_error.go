package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ProcessMsgError tests a process msg returning an error
func ProcessMsgError() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  100,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()
	return &tests.ControllerSpecTest{
		Name: "process msg error ",
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
		ExpectedError: "could not process msg: proposal invalid: proposal not justified: change round has not quorum",
	}
}
