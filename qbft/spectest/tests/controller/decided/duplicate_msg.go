package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// DuplicateMsg tests a duplicate decided msg processing
func DuplicateMsg() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide duplicate msg",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, 10),
					testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, 10),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt:               1,
					DecidedVal:               []byte{1, 2, 3, 4},
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{0, 10},
				},
				ControllerPostRoot: "6ec856c53c4febbeeb0d816b81a04425f5a7bdf107c7cf3d28a519c3fee6ce6e",
			},
		},
	}
}
