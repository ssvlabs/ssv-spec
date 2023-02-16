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
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingProposalMessage(ks.Shares[1], 1),
				},
				ControllerPostRoot: "ede93c772804aff585990dc09b29d841ce55e024831ea05674960b126922a0af",
			},
		},
	}
}
