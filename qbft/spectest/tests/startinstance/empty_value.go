package startinstance

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
)

// EmptyValue tests a starting an instance for an empty value (not passing value check)
func EmptyValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance empty value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{},
				ControllerPostRoot: "8d4a1b5011b185f3b657b7b9e55c82940768031a5f858a623f529d068f1fd28b",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
