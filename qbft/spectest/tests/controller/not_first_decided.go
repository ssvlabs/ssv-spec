package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotFirstDecided tests a process msg after which not first time decided
func NotFirstDecided() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	msgs := testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier, qbft.FirstHeight, testingutils.Testing4SharesSet())
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	}).Encode()

	msgs = append(msgs, &types.Message{
		ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
		Data: signMsgEncoded,
	})

	return &tests.ControllerSpecTest{
		Name: "not first decided",
		RunInstanceData: []struct {
			InputValue    []byte
			InputMessages []*types.Message
			Decided       bool
			DecidedVal    []byte
			DecidedCnt    uint
			SavedDecided  *qbft.SignedMessage
		}{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				Decided:       true,
				DecidedVal:    []byte{1, 2, 3, 4},
				DecidedCnt:    1,
			},
		},
	}
}
