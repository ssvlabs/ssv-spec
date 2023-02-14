package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// FullDecided tests process msg and first time deciding
func FullDecided() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "full decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingProposalMessage(ks.Shares[1], 1),

					testingutils.TestingPrepareMessage(ks.Shares[1], 1),
					testingutils.TestingPrepareMessage(ks.Shares[2], 2),
					testingutils.TestingPrepareMessage(ks.Shares[3], 3),

					testingutils.TestingCommitMessage(ks.Shares[1], 1),
					testingutils.TestingCommitMessage(ks.Shares[2], 2),
					testingutils.TestingCommitMessage(ks.Shares[3], 3),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: []byte{1, 2, 3, 4},
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
					),
				},
				ControllerPostRoot: "f82a7925fa41a67b245d6f97b13c1d272632ac4efe0380847ac8c9378f5bb04b",
			},
		},
	}
}
