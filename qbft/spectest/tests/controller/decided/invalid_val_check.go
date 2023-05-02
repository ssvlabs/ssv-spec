package decided

import (
	"crypto/sha256"

	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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

func invalidValCheckDataStateComparison() *qbftcomparable.StateComparison {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)
	_ = contr.StartNewInstance([]byte{1, 2, 3, 4})

	ks := testingutils.Testing4SharesSet()

	state := testingutils.BaseInstance().State
	state.Decided = true
	state.DecidedValue = testingutils.TestingInvalidValueCheck
	state.CommitContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
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
	}}

	contr.StoredInstances[0].State = state

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
