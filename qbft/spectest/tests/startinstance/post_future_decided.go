package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// PostFutureDecided tests starting a new instance after deciding with future decided msg
func PostFutureDecided() tests.SpecTest {
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
				ControllerPostRoot: "589b0c0352f1c22875246f2e66530d5fda62f646434b250ade128c61c16f47bd",
			},
			{
				InputValue: []byte{1, 2, 3, 4},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal:               testingutils.TestingQBFTFullData,
					DecidedCnt:               0,
					CalledSyncDecidedByRange: true,
					DecidedByRangeValues:     [2]qbft.Height{qbft.FirstHeight, 10},
				},

				ControllerPostRoot: "589b0c0352f1c22875246f2e66530d5fda62f646434b250ade128c61c16f47bd",
			},
		},
		ExpectedError: "invalid instance height",
	}
}
