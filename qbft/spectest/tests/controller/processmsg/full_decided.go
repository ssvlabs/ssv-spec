package processmsg

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	qbftcomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// FullDecided tests process msg and first time deciding
func FullDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := fullDecidedStateComparison()
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

func fullDecidedStateComparison() *qbftcomparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.ExpectedDecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData, testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: testingutils.TestingQBFTFullData,
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
	qbftcomparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
