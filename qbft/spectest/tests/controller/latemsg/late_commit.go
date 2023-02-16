package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateCommit tests process late commit msg for an instance which just decided
func LateCommit() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()

	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.DefaultIdentifier, qbft.FirstHeight, ks)
	msgs = append(msgs, testingutils.TestingCommitMessage(ks.Shares[4], types.OperatorID(4)))

	return &tests.ControllerSpecTest{
		Name: "late commit",
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
				ControllerPostRoot: "61a03b8249dac35278398d36018b6037a8c92868bbc60a2ba4922dffccb30c39",
			},
		},
	}
}
