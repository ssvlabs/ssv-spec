package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnknownSigner tests a single prepare received with an unknown signer
func UnknownSigner() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(5)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare unknown signer",
		Pre:           pre,
		PostRoot:      "470d1a88e97b20eafb08ad9682c10642de27515fff7a8ef3c2d2e97953432357",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: signers not in committee",
	}
}
