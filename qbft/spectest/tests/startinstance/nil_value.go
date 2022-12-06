package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// NilValue tests a starting an instance for a nil value (not passing value check)
func NilValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance nil value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         nil,
				ControllerPostRoot: "475fd29d6449d161b9d2925b73023dce8c28f0fb2faedaeb2f8b8214de08ac69",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
