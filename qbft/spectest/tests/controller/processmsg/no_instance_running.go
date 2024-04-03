package processmsg

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
						[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
						[]types.OperatorID{1, 2, 3},
						50,
					),
					testingutils.TestingProposalMessageWithHeight(ks.NetworkKeys[1], 1, 2),
				},

				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
				},
				ControllerPostRoot: "e8da00ea07e1e5098026373c51e38a681215e12ca4bdeb1f1efbb9d4f3325a92",
			},
		},
		ExpectedError: "instance not found",
	}
}
