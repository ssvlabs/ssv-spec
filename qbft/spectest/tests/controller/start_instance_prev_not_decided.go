package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// StartInstancePreviousNotDecided tests starting an instance when the previous one not decided
func StartInstancePreviousNotDecided() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance prev not decided",
		RunInstanceData: []struct {
			InputValue    []byte
			InputMessages []*qbft.SignedMessage
			Decided       bool
			DecidedVal    []byte
			DecidedCnt    uint
			SavedDecided  *qbft.SignedMessage
		}{
			{
				InputValue: []byte{1, 2, 3, 4},
				Decided:    false,
				DecidedVal: nil,
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				Decided:    false,
				DecidedVal: nil,
			},
		},
		ExpectedError: "can't start new QBFT instance: previous instance hasn't Decided",
	}
}
