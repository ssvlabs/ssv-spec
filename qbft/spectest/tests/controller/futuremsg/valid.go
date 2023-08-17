package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidMsg tests future msg valid msg. This is a valid msg that is not yet ready to be processed.
func ValidMsg() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := []byte{1, 2, 3, 4}

	msg := testingutils.TestingPrepareMessageWithParams(
		ks.Shares[3], 3, 3, 10, identifier[:], testingutils.TestingQBFTRootData)

	// create base controller
	contr := createBaseController()

	return &tests.ControllerSpecTest{
		Name: "future valid msg",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:          []byte{1, 2, 3, 4},
				InputMessages:       []*qbft.SignedMessage{msg},
				ControllerPostState: contr,
			},
		},
		ExpectedError: "future msg from height, could not process",
	}
}

func createBaseController() *qbft.Controller {
	id := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		id[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config)
	return contr
}
