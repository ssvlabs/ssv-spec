package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// StartInstanceInvalidValue tests a starting an instance for an invalid value (not passing value check)
func StartInstanceInvalidValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []struct {
			InputValue    []byte
			InputMessages []*qbft.SignedMessage
			Decided       bool
			DecidedVal    []byte
			DecidedCnt    uint
			SavedDecided  *qbft.SignedMessage
		}{
			{
				InputValue: []byte{1, 2, 3, 3},
				Decided:    false,
				DecidedVal: nil,
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
