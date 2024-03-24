package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
// Instance starts but with an empty value and no proposal gets broadcasted
func InvalidValue() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         testingutils.TestingInvalidValueCheck,
				ControllerPostRoot: "baf3ccea443a6c639b76dccf2d9c4fb5e48318473797de9b55e4d8de48fccc6b",
			},
		},
	}
}
