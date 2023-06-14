package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
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

// UnknownSignerStateComparison returns the expected state comparison for UnknownSigner test.
// The controller is initialized with 4 shares and no messages in its container since the given msg is invalid.
func unknownSignerStateComparison() *comparable.StateComparison {
	identifier := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			Share:   testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:      identifier[:],
			Round:   qbft.FirstRound,
			Height:  qbft.FirstHeight,
			Decided: false,
		},
	}
	comparable.InitContainers(instance)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &comparable.StateComparison{ExpectedState: contr}
}
