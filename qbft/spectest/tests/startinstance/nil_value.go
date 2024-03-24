package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// NilValue tests a starting an instance for a nil value (not passing value check)
// Instance starts but with an empty value and no proposal gets broadcasted
func NilValue() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance nil value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         nil,
				ControllerPostRoot: "baf3ccea443a6c639b76dccf2d9c4fb5e48318473797de9b55e4d8de48fccc6b",
			},
		},
	}
}
