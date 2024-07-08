package decided

import (
	"crypto/rsa"
	"crypto/sha256"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// InvalidValCheckData tests a decided message with invalid decided data (but should pass as it's decided)
func InvalidValCheckData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := invalidValCheckDataStateComparison()

	return &tests.ControllerSpecTest{
		Name: "decide invalid value (should pass)",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingCommitMultiSignerMessageWithParams(
						[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
						[]types.OperatorID{1, 2, 3},
						qbft.FirstRound,
						qbft.FirstHeight,
						testingutils.TestingIdentifier,
						sha256.Sum256(testingutils.TestingInvalidValueCheck),
						testingutils.TestingInvalidValueCheck,
					),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingInvalidValueCheck,
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
	}
}

func invalidValCheckDataStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMultiSignerMessageWithParams(
			[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
			[]types.OperatorID{1, 2, 3},
			qbft.FirstRound,
			qbft.FirstHeight,
			testingutils.TestingIdentifier,
			sha256.Sum256(testingutils.TestingInvalidValueCheck),
			testingutils.TestingInvalidValueCheck,
		),
	}

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			CommitteeMember: testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
			ID:              testingutils.TestingIdentifier,
			Decided:         true,
			DecidedValue:    testingutils.TestingInvalidValueCheck,
			Round:           qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &comparable.StateComparison{ExpectedState: contr}
}
