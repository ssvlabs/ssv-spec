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
		PostRoot:      "01489f7af13579b66ce3da156d4d10208c85a10365380f04e7b8d82d0a9679ce",
		InputMessages: []*types.SignedSSVMessage{msg},
		ExpectedError: "invalid signed message: H(data) != root",
	}
}
