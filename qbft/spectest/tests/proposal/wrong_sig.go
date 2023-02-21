package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSignature tests a single proposal received with a wrong signature
func WrongSignature() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[2], types.OperatorID(1)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong proposal sig",
		Pre:           pre,
		PostRoot:      "e9a7c1b25a014f281d084637bdc50ec144f1262b5cc87eeb6b4493d27d10a69b",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg signature invalid: failed to verify signature",
	}
}
