package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidIdentifier tests a process msg with the wrong identifier
func InvalidIdentifier() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessage(ks.Shares[1], 1)

	return &tests.ControllerSpecTest{
		Name: "invalid identifier",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*qbft.SignedMessage{
					msg,
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
		},
		ExpectedError: "invalid msg: message doesn't belong to Identifier",
	}
}
