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
				ControllerPostRoot: "475fd29d6449d161b9d2925b73023dce8c28f0fb2faedaeb2f8b8214de08ac69",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
