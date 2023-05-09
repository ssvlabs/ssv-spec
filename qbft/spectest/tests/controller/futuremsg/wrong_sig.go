package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSig tests future msg with invalid sig
func WrongSig() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := wrongSigStateComparison()
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)

	return &ControllerSyncSpecTest{
		Name: "future msg wrong sig",
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 2, 3, 10,
				identifier[:], testingutils.TestingQBFTRootData),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   sc.Root(),
		ControllerPostState:  sc.ExpectedState,
		ExpectedError:        "invalid future msg: msg signature invalid: failed to verify signature",
	}
}

func wrongSigStateComparison() *qbftcomparable.StateComparison {
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			Share: testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:    identifier[:],
			Round: qbft.FirstRound,
		},
	}
	qbftcomparable.SetSignedMessages(instance, []*qbft.SignedMessage{})
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
