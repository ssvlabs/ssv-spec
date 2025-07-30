package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
func InvalidValue() tests.SpecTest {
	return tests.NewControllerSpecTest(
		"start instance invalid value",
		testdoc.StartInstanceInvalidValueDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingInvalidValueCheck,
			},
		},
		nil,
		"value invalid: invalid value",
		nil,
		nil,
	)
}
