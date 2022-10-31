package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreviousDecided tests starting an instance when the previous one decided
func PreviousDecided() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "start instance prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				InputMessages:      testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], qbft.FirstHeight, testingutils.Testing4SharesSet()),
				DecidedVal:         inputData.Source,
				DecidedCnt:         1,
				ControllerPostRoot: "d5d4696d29f1359a0f55292ba42dfd922993408529aa86926243df2221554c11",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "ef4b84dc6704519af8f6c4a510a2d9d0a44ce52155f6508635dacbd34324b32e",
			},
		},
	}
}
