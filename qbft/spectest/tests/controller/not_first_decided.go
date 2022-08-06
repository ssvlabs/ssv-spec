package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotFirstDecided tests a process msg after which not first time decided
func NotFirstDecided() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	msgs := testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, testingutils.Testing4SharesSet())
	msgs = append(msgs, testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[4], 4, &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: identifier[:],
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	}))

	return &tests.ControllerSpecTest{
		Name: "not first decided",
		RunInstanceData: []struct {
			InputValue    []byte
			InputMessages []*qbft.SignedMessage
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
