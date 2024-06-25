package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests a single proposal received with an unknown signer
func UnknownSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[2], types.OperatorID(5)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "unknown proposal signer",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: signer not in committee",
	}
}
