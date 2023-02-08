package startinstance

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
)

// PreviousNotDecided tests starting an instance when the previous one not decided
func PreviousNotDecided() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance prev not decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
		},
		ExpectedError: "can't start new QBFT instance: previous instance hasn't Decided",
	}
}
