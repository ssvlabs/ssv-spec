package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownSigner tests a single proposal received with an unknown signer
func UnknownSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[2], types.OperatorID(5)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "unknown proposal signer",
		Pre:           pre,
		InputMessages: msgs,
		ExpectedError: "invalid signed message: signer not in committee",
	}
}
