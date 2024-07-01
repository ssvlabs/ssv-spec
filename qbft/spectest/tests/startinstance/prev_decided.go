package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// PreviousDecided tests starting an instance when the previous one decided
func PreviousDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "start instance prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
					testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks),
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
				},
				ControllerPostRoot:  previousDecidedStateComparison(qbft.FirstHeight, true).Root(),
				ControllerPostState: previousDecidedStateComparison(qbft.FirstHeight, true).ExpectedState,
			},
			{
				InputValue:          []byte{1, 2, 3, 4},
				ControllerPostRoot:  previousDecidedStateComparison(1, false).Root(),
				ControllerPostState: previousDecidedStateComparison(1, false).ExpectedState,
			},
		},
	}
}

func previousDecidedStateComparison(height qbft.Height, decidedState bool) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData, testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	for i := 0; i <= int(height); i++ {
		contr.Height = qbft.Height(i)

		instance := &qbft.Instance{
			StartValue: []byte{1, 2, 3, 4},
			State: &qbft.State{
				CommitteeMember: testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
				ID:              testingutils.TestingIdentifier,
				Round:           qbft.FirstRound,
				Height:          qbft.Height(i),
			},
		}

		// last height
		if !decidedState && qbft.Height(i) == height {
			comparable.SetSignedMessages(instance, []*types.SignedSSVMessage{})
			contr.StoredInstances = append([]*qbft.Instance{instance}, contr.StoredInstances...)
			break
		}

		instance.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), qbft.FirstRound, qbft.Height(i), testingutils.TestingQBFTRootData, nil, nil))
		instance.State.LastPreparedRound = qbft.FirstRound
		instance.State.LastPreparedValue = testingutils.TestingQBFTFullData
		instance.State.Decided = true
		instance.State.DecidedValue = testingutils.TestingQBFTFullData
		if qbft.Height(i) != height {
			instance.ForceStop()
		}

		comparable.SetSignedMessages(instance, msgs)
		contr.StoredInstances = append([]*qbft.Instance{instance}, contr.StoredInstances...)
	}

	return &comparable.StateComparison{ExpectedState: contr}
}
