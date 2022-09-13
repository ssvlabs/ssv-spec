package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreviousDecided tests starting an instance when the previous one decided
func PreviousDecided() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "start instance prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				InputMessages:      testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, testingutils.Testing4SharesSet()),
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         1,
				ControllerPostRoot: "aa402d7487719b17dde352e2ac602ba2c7d895e615ab12cd93d816f6c4fa0967",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "ef4b84dc6704519af8f6c4a510a2d9d0a44ce52155f6508635dacbd34324b32e",
			},
		},
	}
}
