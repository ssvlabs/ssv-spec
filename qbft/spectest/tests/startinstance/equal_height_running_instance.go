package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// EqualHeightRunningInstance tests starting an instance for height equal to a running instance
func EqualHeightRunningInstance() tests.SpecTest {
	height := qbft.FirstHeight

	return tests.NewControllerSpecTest(
		"start instance equal height running instance",
		testdoc.StartInstanceEqualHeightRunningInstanceDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &height,
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				Height:     &height,
			},
		},
		nil,
		"instance already running",
		nil,
	)
}
