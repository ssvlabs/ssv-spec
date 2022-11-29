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
				ControllerPostRoot: "2e8a664c6fa643b691d0b8d56d9819f3c634f6fdb5990d869b4f08c3a1917a47",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
