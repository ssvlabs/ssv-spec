package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
// Instance starts but with an empty value and no proposal gets broadcasted
func InvalidValue() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingInvalidValueCheck,
			},
		},
	}
}
