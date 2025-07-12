package startinstance

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PostFutureDecided tests starting a new instance after deciding with future decided msg
func PostFutureDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return tests.NewControllerSpecTest(
		"start instance post future decided",
		"Test starting a new QBFT instance after a future decided message, then attempting to start with a past height.",
		[]*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
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
				InputValue: []byte{1, 2, 3, 4},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 0,
				},
			},
		},
		nil,
		"attempting to start an instance with a past height",
		nil,
	)
}
