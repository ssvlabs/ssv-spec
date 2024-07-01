package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// LowerHeight tests starting an instance a height lower than previous height
func LowerHeight() tests.SpecTest {
	height1 := qbft.Height(1)
	firstHeight := qbft.FirstHeight

	return &tests.ControllerSpecTest{
		Name: "start instance lower height",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &height1,
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &firstHeight,
			},
		},
		ExpectedError: "attempting to start an instance with a past height",
	}
}
