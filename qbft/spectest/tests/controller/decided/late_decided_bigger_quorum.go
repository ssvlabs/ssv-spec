package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateDecidedBiggerQuorum tests processing a decided msg for a just decided instance (with a bigger quorum)
func LateDecidedBiggerQuorum() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)
	msgs = append(msgs, testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3], ks.Shares[4]}, []types.OperatorID{1, 2, 3, 4}))
	return &tests.ControllerSpecTest{
		Name: "decide late decided bigger quorum",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt:         1,
					DecidedVal:         testingutils.TestingQBFTFullData,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}),
				},
				ControllerPostRoot: "b61697f2ff067f353262fec6946aaf5c0ee94fd56e98128757fe737af83b754d",
			},
		},
	}
}
