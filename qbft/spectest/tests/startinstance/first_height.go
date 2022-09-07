package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// FirstHeight tests a starting the first instance
func FirstHeight() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance first height",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				DecidedVal:         nil,
				ControllerPostRoot: "f28cfa54f7993a21ebe3a46fde514e387ad4ff9a3be197b656c4c6e8dbb124de",
			},
		},
	}
}
