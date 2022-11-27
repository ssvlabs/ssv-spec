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
				ControllerPostRoot: "2e8a664c6fa643b691d0b8d56d9819f3c634f6fdb5990d869b4f08c3a1917a47",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
