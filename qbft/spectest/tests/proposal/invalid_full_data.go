package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidFullData tests signed proposal with an invalid full data field (H(full data) != root)
func InvalidFullData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))
	msg.FullData = nil

	return &tests.MsgProcessingSpecTest{
		Name:          "invalid full data",
		Pre:           pre,
		PostRoot:      "620ad2417e47411537db8df9d4a072327e3c3efc391c3162867f30d5bf9af52c",
		InputMessages: []*types.SignedSSVMessage{msg},
		ExpectedError: "invalid signed message: H(data) != root",
	}
}
