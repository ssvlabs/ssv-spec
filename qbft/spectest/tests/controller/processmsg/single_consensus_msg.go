package processmsg

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SingleConsensusMsg tests process msg of a single msg
func SingleConsensusMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return tests.NewControllerSpecTest(
		"single consensus msg",
		"Test processing a single consensus message.",
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),
				},
			},
		},
		nil,
		"",
		nil,
	)
}
