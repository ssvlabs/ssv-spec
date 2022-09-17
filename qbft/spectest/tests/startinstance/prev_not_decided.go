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
				ControllerPostRoot: "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
			},
		},
		ExpectedError: "can't start new QBFT instance: previous instance hasn't Decided",
	}
}
