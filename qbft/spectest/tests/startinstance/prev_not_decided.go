package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	qbftcomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// PreviousNotDecided tests starting an instance when the previous one not decided
func PreviousNotDecided() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance prev not decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:          []byte{1, 2, 3, 4},
				ControllerPostRoot:  previousNotDecided1SC().Root(),
				ControllerPostState: previousNotDecided1SC().ExpectedState,
			},
			{
				InputValue:          []byte{1, 2, 3, 4},
				ControllerPostRoot:  previousNotDecided2SC().Root(),
				ControllerPostState: previousNotDecided2SC().ExpectedState,
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
	instance := &qbft.Instance{
		State: &qbft.State{
			Share:             testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:                identifier,
			Round:             qbft.FirstRound,
			Height:            qbft.FirstHeight,
			LastPreparedRound: qbft.NoRound,
		},
		StartValue: []byte{1, 2, 3, 4},
	}
	qbftcomparable.SetMessages(instance, []*types.SSVMessage{})
	contr.StoredInstances = append(contr.StoredInstances, instance)
	return &qbftcomparable.StateComparison{ExpectedState: contr}
}

func previousNotDecided2SC() *qbftcomparable.StateComparison {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)
	instance1 := &qbft.Instance{
		State: &qbft.State{
			Share:             testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:                identifier,
			Round:             qbft.FirstRound,
			Height:            qbft.FirstHeight,
			LastPreparedRound: qbft.NoRound,
		},
		StartValue: []byte{1, 2, 3, 4},
	}
	qbftcomparable.SetMessages(instance1, []*types.SSVMessage{})

	instance2 := &qbft.Instance{
		State: &qbft.State{
			Share:             testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:                identifier,
			Round:             qbft.FirstRound,
			Height:            1,
			LastPreparedRound: qbft.NoRound,
		},
		StartValue: []byte{1, 2, 3, 4},
	}
	qbftcomparable.SetMessages(instance2, []*types.SSVMessage{})
	contr.StoredInstances = []*qbft.Instance{instance2, instance1}
	contr.Height = 1
	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
