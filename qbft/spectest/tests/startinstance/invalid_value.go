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
				ControllerPostRoot: "750901b7cf8aed577cf2b2ed23b500c76a77fc2e144e23f3f1dfe6cd8876e3af",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
