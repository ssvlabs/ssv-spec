package futuremsg

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidMsg tests future msg valid msg. This is a valid msg that is not yet ready to be processed.
func ValidMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := testingutils.TestingIdentifier

	msg := testingutils.TestingPrepareMessageWithParams(
		ks.OperatorKeys[3], 3, 3, 10, identifier[:], testingutils.TestingQBFTRootData,
	)

	// create base controller
	contr := createBaseController()

	return tests.NewControllerSpecTest(
		"future valid msg",
		"Test future message that is valid but not yet ready to be processed, expecting error.",
		[]*tests.RunInstanceData{
			{
				InputValue:          []byte{1, 2, 3, 4},
				InputMessages:       []*types.SignedSSVMessage{msg},
				ControllerPostState: contr,
			},
		},
		nil,
		"future msg from height, could not process",
		nil,
	)
}

func createBaseController() *qbft.Controller {
	id := testingutils.TestingIdentifier
	ks := testingutils.Testing4SharesSet()
	config := testingutils.TestingConfig(ks)
	contr := testingutils.NewTestingQBFTController(
		id[:],
		testingutils.TestingCommitteeMember(ks),
		config,
		testingutils.TestingOperatorSigner(ks),
	)
	return contr
}
