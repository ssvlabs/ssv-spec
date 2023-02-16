package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateProposal tests process late proposal msg for an instance which just decided
func LateProposal() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.DefaultIdentifier, qbft.FirstHeight, ks)
	msgs = append(msgs, testingutils.TestingProposalMessage(ks.Shares[1], 1))

	return &tests.ControllerSpecTest{
		Name: "late proposal",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
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
		ExpectedError: "could not process msg: invalid signed message: proposal is not valid with current state",
	}
}
