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
				ControllerPostRoot: "f28cfa54f7993a21ebe3a46fde514e387ad4ff9a3be197b656c4c6e8dbb124de",
			},
			{
				InputValue: []byte{1, 2, 3, 4},
			},
		},
		ExpectedError: "can't start new QBFT instance: previous instance hasn't Decided",
	}
}
