package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// EmptyValue tests a starting an instance for an empty value (not passing value check)
func EmptyValue() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance empty value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{},
			},
		},
		ExpectedError: "value invalid: invalid value",
	}
}
