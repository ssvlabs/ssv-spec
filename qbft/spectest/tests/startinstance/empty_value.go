package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
)

// EmptyValue tests a starting an instance for an empty value (not passing value check)
func EmptyValue() tests.SpecTest {
	return tests.NewControllerSpecTest(
		"start instance empty value",
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{},
			},
		},
		nil,
		"value invalid: invalid value",
		nil,
	)
}
