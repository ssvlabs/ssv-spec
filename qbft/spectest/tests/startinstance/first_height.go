package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// FirstHeight tests a starting the first instance
func FirstHeight() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance first height",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				DecidedVal:         nil,
				ControllerPostRoot: "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
			},
		},
	}
}
