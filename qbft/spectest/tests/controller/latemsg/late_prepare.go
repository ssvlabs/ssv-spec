package latemsg

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// LatePrepare tests process late prepare msg for an instance which just decided
func LatePrepare() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := latePrepareStateComparison()

	msgs := testingutils.DecidingMsgsForHeightWithRoot(
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks,
	)
	msgs = append(msgs, testingutils.TestingPrepareMessage(ks.OperatorKeys[4], 4))

	test := tests.NewControllerSpecTest(
		"late prepare",
		testdoc.ControllerLateMsgLatePrepareDoc,
		[]*tests.RunInstanceData{
			{
				InputValue:    testingutils.TestingQBFTFullData,
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
					),
				},

				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
		nil,
		types.SkipConsensusMessageAsInstanceIsDecidedErrorCode,
		nil,
		ks,
	)

	return test
}

// LatePrepareStateComparison returns the expected state comparison for LatePrepare test.
// The controller is initialized with 4 shares and all expected messages in its container from 3 nodes,
// in addition to the late prepare msg from the 4th node.
// The instance is decided.
func latePrepareStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.ExpectedDecidingMsgsForHeightWithRoot(
		testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData,
		testingutils.TestingIdentifier,
		qbft.FirstHeight,
		ks,
	)
	// append late prepare msg
	msgs = append(msgs, testingutils.TestingPrepareMessage(ks.OperatorKeys[4], types.OperatorID(4)))

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
		testingutils.TestingOperatorSigner(ks),
	)

	instance := &qbft.Instance{
		StartValue: testingutils.TestingQBFTFullData,
		State: &qbft.State{
			CommitteeMember: testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
			ID:              testingutils.TestingIdentifier,
			ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(
				testingutils.TestingProposalMessage(
					ks.OperatorKeys[1],
					types.OperatorID(1),
				),
			),
			LastPreparedRound: qbft.FirstRound,
			LastPreparedValue: testingutils.TestingQBFTFullData,
			Decided:           true,
			DecidedValue:      testingutils.TestingQBFTFullData,
			Round:             qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &comparable.StateComparison{ExpectedState: contr}
}
