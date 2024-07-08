package processmsg

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoInstanceRunning tests a process msg for height in which there is no running instance
func NoInstanceRunning() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "no instance running",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						50,
					),
					testingutils.TestingProposalMessageWithHeight(ks.OperatorKeys[1], 1, 2),
				},

				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
				},
			},
		},
		ExpectedError: "instance not found",
	}
}
