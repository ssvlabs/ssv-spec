package startinstance

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	qbftcomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// PreviousNotDecided tests starting an instance when the previous one not decided
func PreviousNotDecided() tests.SpecTest {
	return tests.NewControllerSpecTest(
		"start instance prev not decided",
		testdoc.StartInstancePrevNotDecidedDoc,
		[]*tests.RunInstanceData{
			{
				InputValue:          testingutils.TestingQBFTFullData,
				ControllerPostRoot:  previousNotDecided1SC().Root(),
				ControllerPostState: previousNotDecided1SC().ExpectedState,
			},
			{
				InputValue:          testingutils.TestingQBFTFullData,
				ControllerPostRoot:  previousNotDecided2SC().Root(),
				ControllerPostState: previousNotDecided2SC().ExpectedState,
			},
		},
		nil,
		0,
		nil,
		nil,
	)
}

func previousNotDecided1SC() *qbftcomparable.StateComparison {
	identifier := testingutils.TestingIdentifier
	ks := testingutils.Testing4SharesSet()
	config := testingutils.TestingConfig(ks)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingCommitteeMember(ks),
		config,
		testingutils.TestingOperatorSigner(ks),
	)
	instance := &qbft.Instance{
		State: &qbft.State{
			CommitteeMember:   testingutils.TestingCommitteeMember(ks),
			ID:                identifier,
			Round:             qbft.FirstRound,
			Height:            qbft.FirstHeight,
			LastPreparedRound: qbft.NoRound,
		},
		StartValue: testingutils.TestingQBFTFullData,
	}
	qbftcomparable.SetMessages(instance, []*types.SignedSSVMessage{})
	contr.StoredInstances = append(contr.StoredInstances, instance)
	return &qbftcomparable.StateComparison{ExpectedState: contr}
}

func previousNotDecided2SC() *qbftcomparable.StateComparison {
	identifier := testingutils.TestingIdentifier
	ks := testingutils.Testing4SharesSet()
	config := testingutils.TestingConfig(ks)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingCommitteeMember(ks),
		config,
		testingutils.TestingOperatorSigner(ks),
	)
	instance1 := &qbft.Instance{
		State: &qbft.State{
			CommitteeMember:   testingutils.TestingCommitteeMember(ks),
			ID:                identifier,
			Round:             qbft.FirstRound,
			Height:            qbft.FirstHeight,
			LastPreparedRound: qbft.NoRound,
		},
		StartValue: testingutils.TestingQBFTFullData,
	}
	qbftcomparable.SetMessages(instance1, []*types.SignedSSVMessage{})
	instance1.ForceStop()

	instance2 := &qbft.Instance{
		State: &qbft.State{
			CommitteeMember:   testingutils.TestingCommitteeMember(ks),
			ID:                identifier,
			Round:             qbft.FirstRound,
			Height:            1,
			LastPreparedRound: qbft.NoRound,
		},
		StartValue: testingutils.TestingQBFTFullData,
	}
	qbftcomparable.SetMessages(instance2, []*types.SignedSSVMessage{})
	contr.StoredInstances = []*qbft.Instance{instance2, instance1}
	contr.Height = 1
	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
