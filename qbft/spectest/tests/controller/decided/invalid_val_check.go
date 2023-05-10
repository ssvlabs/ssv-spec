package decided

import (
	"crypto/sha256"

	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
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
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingCommitMultiSignerMessageWithParams(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
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
	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMultiSignerMessageWithParams(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
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
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			Share:        testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:           testingutils.TestingIdentifier,
			Decided:      true,
			DecidedValue: testingutils.TestingInvalidValueCheck,
			Round:        qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &comparable.StateComparison{ExpectedState: contr}
}
