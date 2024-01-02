package latemsg

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// LateProposalPastInstanceNonDuplicate tests a non-duplicated proposal msg for a previously decided instance
func LateProposalPastInstanceNonDuplicate() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	allMsgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, 1, ks)

	msgPerHeight := make(map[qbft.Height][]*qbft.SignedMessage)
	msgPerHeight[qbft.FirstHeight] = allMsgs[0:7]
	msgPerHeight[1] = allMsgs[7:14]

	instanceData := func(height qbft.Height) *tests.RunInstanceData {
		sc := lateProposalPastInstanceStateComparison(height, nil)
		return &tests.RunInstanceData{
			InputValue:    []byte{1, 2, 3, 4},
			InputMessages: msgPerHeight[height],
			ExpectedDecidedState: tests.DecidedState{
				DecidedVal: testingutils.TestingQBFTFullData,
				DecidedCnt: 1,
				BroadcastedDecided: testingutils.TestingCommitMultiSignerMessageWithHeight(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					height,
				),
			},
			ControllerPostRoot:  sc.Root(),
			ControllerPostState: sc.ExpectedState,
		}
	}

	// Late proposal message for first height and round 2
	lateMsg := testingutils.TestingMultiSignerProposalMessageWithParams([]*bls.SecretKey{ks.Shares[1]}, []types.OperatorID{1}, 2, qbft.FirstHeight, testingutils.TestingIdentifier, testingutils.TestingQBFTFullData, testingutils.TestingQBFTRootData)
	sc := lateProposalPastInstanceStateComparison(2, lateMsg)

	return &tests.ControllerSpecTest{
		Name: "late non duplicate proposal past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight),
			instanceData(1),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					lateMsg,
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
		ExpectedError: "could not process msg: instance stopped processing messages",
	}
}
