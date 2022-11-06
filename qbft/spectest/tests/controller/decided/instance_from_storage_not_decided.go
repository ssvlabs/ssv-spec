package decided

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NotDecidedInstanceFromStorage tests a case where not decided instance is stored in storage
func NotDecidedInstanceFromStorage() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	highestStorageInstance := testingutils.BaseInstance()
	highestStorageInstance.State.ID = identifier[:]

	return &tests.ControllerSpecTest{
		Name:                   "storage highest not decided instance",
		HighestStorageInstance: highestStorageInstance,
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         0,
				ControllerPostRoot: "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
			},
		},
		ExpectedError: "can't start new QBFT instance: previous instance hasn't Decided",
	}
}
