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
		PostRoot:      "35d96d47901971b46055a010e136a97c71644d7d8bee368a8c1f29c149d3fa6c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: signers not in committee",
	}
}
