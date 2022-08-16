package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// StartInstancePreviousDecided tests starting an instance when the previous one decided
func StartInstancePreviousDecided() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "start instance  prev decided",
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
				InputMessages: testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, testingutils.Testing4SharesSet()),
				Decided:       true,
				DecidedVal:    []byte{1, 2, 3, 4},
				DecidedCnt:    1,
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				Decided:    false,
				DecidedVal: nil,
			},
		},
	}
}
