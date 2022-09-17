package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
func InvalidValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         testingutils.TestingInvalidValueCheck,
				ControllerPostRoot: "83cf6310fae8c6985653f4727f849f68c47fc31c9b10c7223a0935c97669bbb4",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
