package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// NoInstanceRunning tests a process msg for height in which there is no running instance
func NoInstanceRunning() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "no instance running",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						50,
					),
					testingutils.TestingProposalMessageWithHeight(ks.Shares[1], 1, 2),
				},

				ExpectedDecidedState: tests.DecidedState{
					DecidedVal:               testingutils.TestingQBFTFullData,
					DecidedCnt:               1,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{0, 50},
				},
				ControllerPostRoot: "e73515e6ec3c6766264dd4f51c205c9aa0e0d598ee4f9c1705dd6aa74e9e96e3",
			},
		},
		ExpectedError: "instance not found",
	}
}
