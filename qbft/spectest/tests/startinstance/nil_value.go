package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// NilValue tests a starting an instance for a nil value (not passing value check)
func NilValue() tests.SpecTest {
	return tests.NewControllerSpecTest(
		"start instance nil value",
		"Test starting a new QBFT instance with a nil value, expecting value validation error.",
		[]*tests.RunInstanceData{
			{
				InputValue: nil,
			},
		},
		nil,
		"value invalid: invalid value",
		nil,
	)
}
