package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidValueCheck tests a proposal that doesn't pass value check
func InvalidValueCheck() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithIdentifierAndFullData(
			ks.OperatorKeys[1], types.OperatorID(1), testingutils.TestingIdentifier, testingutils.TestingInvalidValueCheck,
			qbft.FirstHeight),
	}

	return tests.NewMsgProcessingSpecTest(
		"invalid proposal value check",
		"Test proposal that fails the value check validation, expecting validation error.",
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: proposal not justified: proposal fullData invalid: invalid value",
		nil,
	)
}
