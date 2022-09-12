package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
)

// StartInstanceFirstHeight tests a starting the first instance
func StartInstanceFirstHeight() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance first height",
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
				Decided:    false,
				DecidedVal: nil,
			},
		},
	}
}
