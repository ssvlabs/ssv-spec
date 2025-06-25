package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// LowerHeight tests starting an instance a height lower than previous height
func LowerHeight() tests.SpecTest {
	height1 := qbft.Height(1)
	firstHeight := qbft.FirstHeight

	return tests.NewControllerSpecTest(
		"start instance lower height",
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &height1,
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &firstHeight,
			},
		},
		nil,
		"attempting to start an instance with a past height",
		nil,
	)
}
