package decided

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidFullData tests signed decided with an invalid full data field
func InvalidFullData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
		[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
		[]types.OperatorID{1, 2, 3},
		10,
		testingutils.TestingIdentifier,
		[]byte("invalid"),
	)

	return &tests.ControllerSpecTest{
		Name: "decide invalid full data",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					msg,
				},
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
			},
		},
		ExpectedError: "invalid decided msg: H(data) != root",
	}
}
