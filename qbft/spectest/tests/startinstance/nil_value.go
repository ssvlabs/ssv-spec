package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// NilValue tests a starting an instance for a nil value (not passing value check)
func NilValue() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance nil value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         nil,
				ControllerPostRoot: "83cf6310fae8c6985653f4727f849f68c47fc31c9b10c7223a0935c97669bbb4",
			},
		},
		ExpectedError: "can't start new QBFT instance: value invalid: invalid value",
	}
}
