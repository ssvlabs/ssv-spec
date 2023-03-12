package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// FullDecided tests process msg and first time deciding
func FullDecided() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "full decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: testingutils.TestingQBFTFullData,
				InputMessages: testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
					testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks),
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
					),
				},
				ControllerPostRoot: "4c96913e87aa17c9f0d5c1b6b220cbc7a66b7b40ef55f1059f1e1fa9f59c94d9",
			},
		},
	}
}
