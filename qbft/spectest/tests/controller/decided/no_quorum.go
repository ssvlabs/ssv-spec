package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoQuorum tests decided msg with < unique 2f+1 signers
func NoQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	test := tests.NewControllerSpecTest(
		"decide no quorum",
		testdoc.ControllerDecidedNoQuorumDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2]},
						[]types.OperatorID{1, 2},
						qbft.FirstHeight,
					),
				},
			},
		},
		nil,
		// TODO: before merge ask engineering how often they see such message in production
		types.NoProposalForCurrentRoundErrorCode,
		nil,
		ks,
	)

	return test
}
