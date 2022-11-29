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
				ControllerPostRoot: "2e8a664c6fa643b691d0b8d56d9819f3c634f6fdb5990d869b4f08c3a1917a47",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
