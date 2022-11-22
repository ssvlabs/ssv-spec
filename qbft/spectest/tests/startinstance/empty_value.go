package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// EmptyValue tests a starting an instance for an empty value (not passing value check)
func EmptyValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance empty value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{},
				ControllerPostRoot: "83cf6310fae8c6985653f4727f849f68c47fc31c9b10c7223a0935c97669bbb4",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
