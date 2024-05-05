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
		PostRoot:      "eaa7264b5d6f05cfcdec3158fcae4ff58c3de1e7e9e12bd876177a58686994d4",
		InputMessages: []*types.SignedSSVMessage{msg},
		ExpectedError: "invalid signed message: H(data) != root",
	}
}
