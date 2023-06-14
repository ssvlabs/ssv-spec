package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	qbftcomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// PreviousNotDecided tests starting an instance when the previous one not decided
func PreviousNotDecided() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance prev not decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: previousNotDecided1SC().Root(),
			},
			{
				InputValue:         []byte{1, 2, 3, 4},
				ControllerPostRoot: previousNotDecided2SC().Root(),
			},
		},
	}
}

func previousNotDecided1SC() *qbftcomparable.StateComparison {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)
	contr.StartNewInstance(qbft.FirstHeight, []byte{1, 2, 3, 4})
	return &qbftcomparable.StateComparison{ExpectedState: contr.StoredInstances[0].State}
}

func previousNotDecided2SC() *qbftcomparable.StateComparison {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)
	contr.StartNewInstance(qbft.FirstHeight, []byte{1, 2, 3, 4})
	contr.StartNewInstance(1, []byte{1, 2, 3, 4})
	return &qbftcomparable.StateComparison{ExpectedState: contr.StoredInstances[1].State}
}
