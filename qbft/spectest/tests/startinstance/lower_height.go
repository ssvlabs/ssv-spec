package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// LowerHeight tests starting an instance a height lower than previous height
func LowerHeight() tests.SpecTest {
	height1 := qbft.Height(1)
	firstHeight := qbft.FirstHeight

	return tests.NewControllerSpecTest(
		"start instance lower height",
		testdoc.StartInstanceLowerHeightDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				Height:     &height1,
			},
			{
				InputValue: testingutils.TestingQBFTFullData,
				Height:     &firstHeight,
			},
		},
		nil,
		"attempting to start an instance with a past height",
		nil,
		nil,
	)
}
