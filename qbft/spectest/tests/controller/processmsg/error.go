package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgError tests a process msg returning an error
func MsgError() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "process msg error",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingProposalMessageWithRound(ks.Shares[1], 1, 100),
				},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
		},
		ExpectedError: "could not process msg: invalid signed message: proposal not justified: change round has no quorum",
	}
}
