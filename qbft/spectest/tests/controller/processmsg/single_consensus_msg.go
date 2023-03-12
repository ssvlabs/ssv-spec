package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SingleConsensusMsg tests process msg of a single msg
func SingleConsensusMsg() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "single consensus msg",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingProposalMessage(ks.Shares[1], 1),
				},
				ControllerPostRoot: "c9258515e169e330c535d38068f2dc6bf3f61e6d36a941ea41a0133435afae22",
			},
		},
	}
}
