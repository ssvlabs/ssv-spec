package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
func InvalidValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 3},
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
