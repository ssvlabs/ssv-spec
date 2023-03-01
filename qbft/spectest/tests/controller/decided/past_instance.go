package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PastInstance tests a decided msg received for past instance
func PastInstance() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide past instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, 100),
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, 80),
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, 90),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt:               3,
					DecidedVal:               testingutils.TestingQBFTFullData,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{0, 100},
				},
				ControllerPostRoot: "d2f7f4bfc09d8695021a3c10657907e6196adda7ff8f06c9b48a368539a2e7bf",
			},
		},
	}
}
