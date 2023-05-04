package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests future msg signed by unknown signer
func UnknownSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := unknownSignerStateComparison()

	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	msg := testingutils.TestingPrepareMessageWithParams(ks.Shares[3], 3, 3, 10,
		identifier[:], testingutils.TestingQBFTRootData)
	msg.Signers = []types.OperatorID{10}

	return &ControllerSyncSpecTest{
		Name: "future msg unknown signer",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   sc.Root(),
		ControllerPostState:  sc.ExpectedState,
		ExpectedError:        "invalid future msg: msg signature invalid: unknown signer",
	}
}

func unknownSignerStateComparison() *qbftcomparable.StateComparison {
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)
	_ = contr.StartNewInstance([]byte{1, 2, 3, 4})

	state := testingutils.BaseInstance().State
	state.ID = identifier[:]
	contr.StoredInstances[0].State = state

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
