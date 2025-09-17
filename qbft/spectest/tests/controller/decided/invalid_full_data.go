package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidFullData tests signed decided with an invalid full data field
func InvalidFullData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessageWithHeight(
		[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
		[]types.OperatorID{1, 2, 3},
		10,
	)
	msg.FullData = []byte("invalid")

	test := tests.NewControllerSpecTest(
		"decide invalid full data",
		testdoc.ControllerDecidedInvalidFullDataDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*types.SignedSSVMessage{
					msg,
				},
			},
		},
		nil,
		"invalid decided msg: H(data) != root",
		nil,
		ks,
	)

	return test
}
