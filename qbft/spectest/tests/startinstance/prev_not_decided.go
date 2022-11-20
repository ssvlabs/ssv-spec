package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// PreviousNotDecided tests starting an instance when the previous one not decided
func PreviousNotDecided() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance prev not decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "3cf8cd7d050943781b25a611cc7b51bc051608f316556a843936b62c276984cb",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "3cf8cd7d050943781b25a611cc7b51bc051608f316556a843936b62c276984cb",
			},
		},
		ExpectedError: "can't start new QBFT instance: previous instance hasn't Decided",
	}
}
