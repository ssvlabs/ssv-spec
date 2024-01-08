package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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

	return &tests.ControllerSpecTest{
		Name: "out or order flow",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    testingutils.TestingQBFTFullData,
				InputMessages: outOfOrderMessages,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
			},
		},
		ExpectedError: "could not process msg: invalid signed message: did not receive proposal for this round",
	}
}
