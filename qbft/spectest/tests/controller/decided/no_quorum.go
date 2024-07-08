package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoQuorum tests decided msg with < unique 2f+1 signers
func NoQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "decide no quorum",
		RunInstanceData: []*tests.RunInstanceData{
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
		// TODO: before merge ask engineering how often they see such message in production
		ExpectedError: "could not process msg: invalid signed message: did not receive proposal for this round",
	}
}
