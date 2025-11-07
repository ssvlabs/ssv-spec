package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
)

// NilValue tests a starting an instance for a nil value (not passing value check)
func NilValue() tests.SpecTest {
	return tests.NewControllerSpecTest(
		"start instance nil value",
		testdoc.StartInstanceNilValueDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: nil,
			},
		},
		nil,
		types.QBFTValueInvalidErrorCode,
		nil,
		nil,
	)
}
