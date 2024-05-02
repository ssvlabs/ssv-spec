package latemsg

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// LateProposalPastInstanceNonDuplicate tests a non-duplicated proposal msg for a previously decided instance
func LateProposalPastInstanceNonDuplicate() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	allMsgsForFirstHeight := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)

	scFirstHeight := lateProposalPastInstanceStateComparison(qbft.FirstHeight, nil)
	instanceData := &tests.RunInstanceData{
		InputValue:    []byte{1, 2, 3, 4},
		InputMessages: allMsgsForFirstHeight,
		ExpectedDecidedState: tests.DecidedState{
			DecidedVal: testingutils.TestingQBFTFullData,
			DecidedCnt: 1,
			BroadcastedDecided: testingutils.TestingCommitMultiSignerMessageWithHeight(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				qbft.FirstHeight,
			),
		},
		ControllerPostRoot:  scFirstHeight.Root(),
		ControllerPostState: scFirstHeight.ExpectedState,
	}

	// Late proposal message for first height and round 2
	lateMsg := testingutils.TestingMultiSignerProposalMessageWithParams([]*bls.SecretKey{ks.Shares[1]}, []types.OperatorID{1}, 2, qbft.FirstHeight, testingutils.TestingIdentifier, testingutils.TestingQBFTFullData, testingutils.TestingQBFTRootData)

	sc := lateProposalPastInstanceStateComparison(1, lateMsg)

	return &tests.ControllerSpecTest{
		Name: "late non duplicate proposal past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData,
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
