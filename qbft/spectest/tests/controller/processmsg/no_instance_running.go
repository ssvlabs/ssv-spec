package processmsg

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoInstanceRunning tests a process msg for lower height in which there is no running instance
func NoInstanceRunning() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	height := qbft.Height(50)
	return &tests.ControllerSpecTest{
		Name: "no instance running",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithHeight(ks.OperatorKeys[1], 1, 0),
				},

				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
					DecidedCnt: 0,
				},
			},
		},
		StartHeight:   &height,
		ExpectedError: "instance not found",
	}
}
