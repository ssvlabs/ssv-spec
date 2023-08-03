package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidValueCheck tests a proposal that doesn't pass value check
func InvalidValueCheck() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithIdentifierAndFullData(
			ks.Shares[1], types.OperatorID(1), []byte{1, 2, 3, 4}, testingutils.TestingInvalidValueCheck,
			qbft.FirstHeight),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "invalid proposal value check",
		Pre:            pre,
		PostRoot:       "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: proposal fullData invalid: invalid value",
	}
}
