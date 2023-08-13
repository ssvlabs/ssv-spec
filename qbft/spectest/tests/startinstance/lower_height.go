package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// LowerHeight tests starting an instance a height lower than previous height
func LowerHeight() tests.SpecTest {
	height1 := qbft.Height(1)
	firstHeight := qbft.FirstHeight

	return &tests.ControllerSpecTest{
		Name: "start instance lower height",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				Height:             &height1,
				ControllerPostRoot: "38a8a319f8deae3fa50442bbe60a87536ace6413a7aebf38094a3129e7ae59ae",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				Height:             &firstHeight,
				ControllerPostRoot: "38a8a319f8deae3fa50442bbe60a87536ace6413a7aebf38094a3129e7ae59ae",
			},
		},
		ExpectedError: "attempting to start an instance with a past height",
	}
}
