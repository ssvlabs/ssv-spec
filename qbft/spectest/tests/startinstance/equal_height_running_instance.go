package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// EqualHeightRunningInstance tests starting an instance for height equal to a running instance
func EqualHeightRunningInstance() tests.SpecTest {
	height := qbft.FirstHeight

	return tests.NewControllerSpecTest(
		"start instance equal height running instance",
		testdoc.StartInstanceEqualHeightRunningInstanceDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				Height:     &height,
			},
			{
				InputValue: testingutils.TestingQBFTFullData,
				Height:     &height,
			},
		},
		nil,
		types.InstanceAlreadyRunningErrorCode,
		nil,
		nil,
	)
}
