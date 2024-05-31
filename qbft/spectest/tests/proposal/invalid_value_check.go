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
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithIdentifierAndFullData(
			ks.OperatorKeys[1], types.OperatorID(1), []byte{1, 2, 3, 4}, testingutils.TestingInvalidValueCheck,
			qbft.FirstHeight),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "invalid proposal value check",
		Pre:            pre,
		PostRoot:       "57e323705826bc5d475ead7f48015a785306be56b33b9fed163fdedf03743754",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: proposal fullData invalid: invalid value",
	}
}
