package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PastInstance tests a decided msg received for past instance
func PastInstance() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide past instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 100),
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 80),
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, 90),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 3,
					DecidedVal: testingutils.TestingQBFTFullData,
				},
			},
		},
	}
}
