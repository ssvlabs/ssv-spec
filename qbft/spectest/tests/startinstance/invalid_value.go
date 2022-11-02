package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
func InvalidValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         &qbft.Data{Root: [32]byte{1, 1, 1, 1}, Source: []byte{1, 1, 1, 1}},
				ControllerPostRoot: "ca1e02c04b7792915149c0e52bbeeb5dc3b244c02a124d43f634473f8d32c577",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
