package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// NilValue tests a starting an instance for a nil value (not passing value check)
func NilValue() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance nil value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: nil,
			},
		},
		ExpectedError: "value invalid: invalid value",
	}
}
