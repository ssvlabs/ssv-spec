package decided

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// CurrentInstance tests a decided msg received for current running instance
func CurrentInstance() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := currentInstanceStateComparison()

	return &tests.ControllerSpecTest{
		Name: "decide current instance",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessage(ks.NetworkKeys[1], 1),

					testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
					testingutils.TestingPrepareMessage(ks.NetworkKeys[2], 2),
					testingutils.TestingPrepareMessage(ks.NetworkKeys[3], 3),

					testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
					testingutils.TestingCommitMessage(ks.NetworkKeys[2], 2),
					testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]}, []types.OperatorID{1, 2, 3}),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingQBFTFullData,
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
	}
}

func currentInstanceStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], types.OperatorID(3)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingCommitMessage(ks.NetworkKeys[2], types.OperatorID(2)),
		testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]}, []types.OperatorID{1, 2, 3}),
	}

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			Share:                           testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:                              testingutils.TestingIdentifier,
			ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1)),
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
