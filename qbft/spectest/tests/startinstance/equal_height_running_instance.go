package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
)

// EqualHeightRunningInstance tests starting an instance for height equal to a running instance
func EqualHeightRunningInstance() tests.SpecTest {
	height := qbft.FirstHeight

	return &tests.ControllerSpecTest{
		Name: "start instance equal height running instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				Height:             &height,
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				Height:             &height,
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
			},
		},
		ExpectedError: "instance already running",
	}
}
