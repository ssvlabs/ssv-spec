package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate decided msg processing
func DuplicateMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	test := tests.NewControllerSpecTest(
		"decide duplicate msg",
		testdoc.ControllerDecidedDuplicateMsgDoc,
		[]*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 10),
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 10),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingQBFTFullData,
				},
			},
		},
		nil,
		"",
		nil,
		ks,
	)

	return test
}
