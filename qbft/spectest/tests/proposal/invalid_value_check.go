package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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

	test := tests.NewMsgProcessingSpecTest(
		"invalid proposal value check",
		testdoc.ProposalInvalidValueCheckDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		types.QBFTValueInvalidErrorCode,
		nil,
		ks,
	)

	return test
}
