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
				InputValue: []byte{1, 2, 3, 4},
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
		"no proposal for round: invalid signed message: did not receive proposal for this round",
		nil,
		ks,
	)

	return test
}
