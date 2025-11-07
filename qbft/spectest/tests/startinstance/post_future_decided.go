package startinstance

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostFutureDecided tests starting a new instance after deciding with future decided msg
func PostFutureDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	test := tests.NewControllerSpecTest(
		"start instance post future decided",
		testdoc.StartInstancePostFutureDecidedDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 10,
					),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
				},
			},
			{
				InputValue: testingutils.TestingQBFTFullData,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 0,
				},
			},
		},
		nil,
		types.StartInstanceErrorCode,
		nil,
		ks,
	)

	return test
}
