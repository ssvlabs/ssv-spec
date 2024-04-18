package decided

import (
	"github.com/herumi/bls-eth-go-binary/bls"
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
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2]},
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
