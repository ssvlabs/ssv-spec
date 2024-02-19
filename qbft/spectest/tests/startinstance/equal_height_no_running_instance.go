package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// EqualHeightNoRunningInstance tests starting an instance for height equal to current without a running instance
func EqualHeightNoRunningInstance() tests.SpecTest {
	height := qbft.Height(2)

	return &tests.ControllerSpecTest{
		Name: "start instance for current height with no running instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &height,
			},
		},
		StartHeight: &height,
	}
}
