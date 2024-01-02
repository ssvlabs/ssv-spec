package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	qbftcomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// OutOfOrderFlow tests a QBFT execution with messages out of order
// One prepare and one commit are received before the proposal causing the instance not to decide
func OutOfOrderFlow() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	correctFlowMessgaes := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)

	// shuffles messages
	outOfOrderMessages := []*qbft.SignedMessage{
		correctFlowMessgaes[3], // prepare 3
		correctFlowMessgaes[6], // commit 3
		correctFlowMessgaes[0], // proposal 1
		correctFlowMessgaes[1], // prepare 1
		correctFlowMessgaes[2], // prepare 2
		correctFlowMessgaes[4], // commit 1
		correctFlowMessgaes[5], // commit 2
	}

	sc := OutOfOrderFlowStateComparison(outOfOrderMessages)

	return &tests.ControllerSpecTest{
		Name: "out or order flow",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    testingutils.TestingQBFTFullData,
				InputMessages: outOfOrderMessages,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
		ExpectedError: "could not process msg: invalid signed message: did not receive proposal for this round",
	}
}

func OutOfOrderFlowStateComparison(msgs []*qbft.SignedMessage) *qbftcomparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: testingutils.TestingQBFTFullData,
		State: &qbft.State{
			Share:                           testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:                              testingutils.TestingIdentifier,
			ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
			LastPreparedRound:               qbft.NoRound,
			LastPreparedValue:               nil,
			Decided:                         false,
			DecidedValue:                    nil,
			Round:                           qbft.FirstRound,
		},
	}
	qbftcomparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
