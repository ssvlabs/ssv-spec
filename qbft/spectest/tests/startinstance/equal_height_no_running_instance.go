package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// EqualHeightNoRunningInstance tests starting an instance for height equal to current without a running instance
func EqualHeightNoRunningInstance() tests.SpecTest {
	height := qbft.Height(2)

	return tests.NewControllerSpecTest(
		"start instance for current height with no running instance",
		testdoc.StartInstanceEqualHeightNoRunningInstanceDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &height,
			},
		},
		nil,
		"",
		&height,
	)
}
