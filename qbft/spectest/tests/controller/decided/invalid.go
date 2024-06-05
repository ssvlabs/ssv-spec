package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Invalid tests decided msg where msg.validate() != nil
func Invalid() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessageWithHeight(
		[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
		[]types.OperatorID{1, 2, 3},
		qbft.FirstHeight,
	)
	msg.OperatorIDs = []types.OperatorID{}
	return &tests.ControllerSpecTest{
		Name: "decide invalid msg",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					msg,
				},
			},
		},
		ExpectedError: "could not process msg: invalid signed message: invalid SignedSSVMessage: no signers",
	}
}
