package processmsg

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// BroadcastedDecided tests broadcasting decided
func BroadcastedDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := broadcastedDecidedStateComparison()
	return &tests.ControllerSpecTest{
		Name: "broadcast decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
					testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks),
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
	}
}

func broadcastedDecidedStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.ExpectedDecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData, testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			CommitteeMember:                 testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
			ID:                              testingutils.TestingIdentifier,
			ProposalAcceptedForCurrentRound: testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))),
			LastPreparedRound:               qbft.FirstRound,
			LastPreparedValue:               testingutils.TestingQBFTFullData,
			Decided:                         true,
			DecidedValue:                    testingutils.TestingQBFTFullData,
			Round:                           qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &comparable.StateComparison{ExpectedState: contr}
}
