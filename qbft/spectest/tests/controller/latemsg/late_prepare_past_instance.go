package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LatePreparePastInstance tests process prepare msg for a previously decided instance
func LatePreparePastInstance() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()

	allMsgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.DefaultIdentifier, 5, ks)

	msgPerHeight := make(map[qbft.Height][]*qbft.SignedMessage)
	msgPerHeight[qbft.FirstHeight] = allMsgs[0:7]
	msgPerHeight[1] = allMsgs[7:14]
	msgPerHeight[2] = allMsgs[14:21]
	msgPerHeight[3] = allMsgs[21:28]
	msgPerHeight[4] = allMsgs[28:35]
	msgPerHeight[5] = allMsgs[35:42]

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue:    []byte{1, 2, 3, 4},
			InputMessages: msgPerHeight[height],
			ExpectedDecidedState: tests.DecidedState{
				BroadcastedDecided: testingutils.TestingCommitMultiSignerMessageWithHeight(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					height,
				),
				DecidedVal: []byte{1, 2, 3, 4},
				DecidedCnt: 1,
			},
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "late prepare past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "f82a7925fa41a67b245d6f97b13c1d272632ac4efe0380847ac8c9378f5bb04b"),
			instanceData(1, "d1b707e4b2251967767d9656dd89cb807460b8dabbfd468772b3c088d89fd71b"),
			instanceData(2, "bff7466de84c53b1b0c39e1f5c9faf3e336622218f583a817b862f55bbf9023d"),
			instanceData(3, "05944813dfd352d5b4ce730c09bbb076ade52689111ce94b201547865fe28844"),
			instanceData(4, "848261610a945d4aa24174fe73471ba2c3b85f1147c9fc5a704ff77c3f1a7bbb"),
			instanceData(5, "bd7d5dc577276a5262d188270dcee321198349aea6eb19e6b6d5446d3bbcd827"),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingPrepareMultiSignerMessageWithHeight(
						[]*bls.SecretKey{ks.Shares[4]},
						[]types.OperatorID{4},
						4,
					),
				},
				ControllerPostRoot: "b6a0be59841bd89b1d7cebfab60a21994b5961861e7a6b5031208175629415ad",
			},
		},
	}
}
