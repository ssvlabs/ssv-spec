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
		PostRoot:       "16577c5ca1fcfd7e43549b13fc730edfd05ac8b0110e451e44d512f034ec6eb3",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: proposal fullData invalid: invalid value",
	}
}
