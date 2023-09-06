package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSignature tests a single proposal received with a wrong signature
func WrongSignature() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[2], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong proposal sig",
		Pre:           pre,
		PostRoot:      "1256adaeaaec7689542238eb7fb94c96bcf430f3b0231be2779755d1fb7f43ab",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg signature invalid: failed to verify signature",
	}
}
