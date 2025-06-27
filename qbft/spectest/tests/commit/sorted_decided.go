package commit

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SortedDecided tests the creation of the decided message that should have sorted signers
func SortedDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return tests.NewControllerSpecTest(
		"sorted decided",
		"Test the creation of the decided message that should have sorted signers",
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),
					testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
					testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
					testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),
					testingutils.TestingCommitMessage(ks.OperatorKeys[4], 4),
					testingutils.TestingCommitMessage(ks.OperatorKeys[2], 2),
					testingutils.TestingCommitMessage(ks.OperatorKeys[3], 3),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage(
						[]*rsa.PrivateKey{ks.OperatorKeys[2], ks.OperatorKeys[3], ks.OperatorKeys[4]},
						[]types.OperatorID{2, 3, 4}),
				},
			},
		},
		nil,
		"",
		nil,
	)
}
