package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
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
		"value invalid: invalid value",
		nil,
		nil,
	)
}
