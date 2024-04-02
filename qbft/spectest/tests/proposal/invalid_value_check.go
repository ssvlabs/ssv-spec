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
		PostRoot:       "35d96d47901971b46055a010e136a97c71644d7d8bee368a8c1f29c149d3fa6c",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: proposal fullData invalid: invalid value",
	}
}
