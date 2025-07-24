package latemsg

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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
				[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
				[]types.OperatorID{1, 2, 3},
				qbft.FirstHeight,
			),
		},
		ControllerPostRoot:  scFirstHeight.Root(),
		ControllerPostState: scFirstHeight.ExpectedState,
	}

	// Late proposal message for first height and round 2
	lateMsg := testingutils.TestingMultiSignerProposalMessageWithParams([]*rsa.PrivateKey{ks.OperatorKeys[1]}, []types.OperatorID{1}, 2, qbft.FirstHeight, testingutils.TestingIdentifier, testingutils.TestingQBFTFullData, testingutils.TestingQBFTRootData)

	sc := lateProposalPastInstanceStateComparison(1, lateMsg)

	test := tests.NewControllerSpecTest(
		"late non duplicate proposal past instance",
		testdoc.ControllerLateMsgLateProposalPastInstanceNonDuplicateDoc,
		[]*tests.RunInstanceData{
			instanceData,
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					lateMsg,
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
		nil,
		"not processing consensus message since instance is already decided",
		nil,
		ks,
	)

	return test
}
