package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownSigner tests a single proposal received with an unknown signer
func UnknownSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[2], types.OperatorID(5)),
	}

	test := tests.NewMsgProcessingSpecTest(
		"unknown proposal signer",
		testdoc.ProposalUnknownSignerDoc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		"invalid signed message: signer not in committee",
		nil,
		ks,
	)

	return test
}
