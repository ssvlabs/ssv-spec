package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PostFutureDecided tests starting a new instance after deciding with future decided msg
func PostFutureDecided() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "start instance post future decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingCommitMultiSignerMessageWithHeight(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, 10,
					),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal:               testingutils.TestingQBFTFullData,
					DecidedCnt:               1,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{qbft.FirstHeight, 10},
				},
				ControllerPostRoot: "6ec856c53c4febbeeb0d816b81a04425f5a7bdf107c7cf3d28a519c3fee6ce6e",
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal:               testingutils.TestingQBFTFullData,
					DecidedCnt:               0,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{qbft.FirstHeight, 10},
				},
				ControllerPostRoot: "9255cf35e771cc91f931f1ed291137c7de4229d3fd2b46a638f1e214bd2a9a04",
			},
		},
	}
}
