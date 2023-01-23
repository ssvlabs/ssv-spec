package startinstance

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
func InvalidValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         testingutils.TestingInvalidValueCheck,
				ControllerPostRoot: "8d4a1b5011b185f3b657b7b9e55c82940768031a5f858a623f529d068f1fd28b",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
