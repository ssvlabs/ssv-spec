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
				ControllerPostRoot: "750901b7cf8aed577cf2b2ed23b500c76a77fc2e144e23f3f1dfe6cd8876e3af",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
