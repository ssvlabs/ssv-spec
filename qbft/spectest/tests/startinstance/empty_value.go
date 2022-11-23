package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// EmptyValue tests a starting an instance for an empty value (not passing value check)
func EmptyValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance empty value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{},
				ControllerPostRoot: "750901b7cf8aed577cf2b2ed23b500c76a77fc2e144e23f3f1dfe6cd8876e3af",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
