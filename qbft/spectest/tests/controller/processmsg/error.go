package processmsg

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgError tests a process msg returning an error
func MsgError() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "process msg error",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], 1, 100),
				},
			},
		},
		ExpectedError: "could not process msg: invalid signed message: proposal not justified: change round has no quorum",
	}
}
