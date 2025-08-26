package futuremsg

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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

	test := tests.NewControllerSpecTest(
		"future valid msg",
		testdoc.ControllerFutureMsgValidDoc,
		[]*tests.RunInstanceData{
			{
				InputValue:    testingutils.TestingQBFTFullData,
				InputMessages:       []*types.SignedSSVMessage{msg},
				ControllerPostState: contr,
			},
		},
		nil,
		"future msg from height, could not process",
		nil,
		ks,
	)

	return test
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
