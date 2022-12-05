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
				ControllerPostRoot: "7b74be21fcdae2e7ed495882d1a499642c15a7f732f210ee84fb40cc97d1ce96",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: "7b74be21fcdae2e7ed495882d1a499642c15a7f732f210ee84fb40cc97d1ce96",
			},
		},
		ExpectedError: "can't start new QBFT instance: previous instance hasn't Decided",
	}
}
