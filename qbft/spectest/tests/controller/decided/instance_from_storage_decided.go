package decided

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InstanceFromStorage tests a case where decided instance is stored in storage
func InstanceFromStorage() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	highestStorageInstance := testingutils.BaseInstance()
	highestStorageInstance.State.ID = identifier[:]
	highestStorageInstance.State.Decided = true

	return &tests.ControllerSpecTest{
		Name:                   "storage highest decided instance",
		HighestStorageInstance: highestStorageInstance,
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         0,
				ControllerPostRoot: "1ae0dafa0b5a04d6b5eafcad54788284d4aafb2075ef50c74103468e2694b49d",
			},
		},
	}
}
