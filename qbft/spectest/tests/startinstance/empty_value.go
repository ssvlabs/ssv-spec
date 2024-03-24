package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// EmptyValue tests a starting an instance for an empty value (not passing value check)
// Instance starts but with an empty value and no proposal gets broadcasted
func EmptyValue() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance empty value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{},
				ControllerPostRoot: "baf3ccea443a6c639b76dccf2d9c4fb5e48318473797de9b55e4d8de48fccc6b",
			},
		},
	}
}
