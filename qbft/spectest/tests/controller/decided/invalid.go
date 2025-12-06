package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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
	test := tests.NewControllerSpecTest(
		"decide invalid msg",
		testdoc.ControllerDecidedInvalidDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*types.SignedSSVMessage{
					msg,
				},
			},
		},
		nil,
		types.NoSignersErrorCode,
		nil,
		ks,
	)

	return test
}
