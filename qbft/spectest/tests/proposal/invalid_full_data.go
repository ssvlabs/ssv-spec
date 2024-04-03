package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidFullData tests signed proposal with an invalid full data field (H(full data) != root)
func InvalidFullData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.NetworkKeys[1], types.OperatorID(1), []byte{1, 2, 3, 4}, nil, qbft.FirstHeight)

	return &tests.MsgProcessingSpecTest{
		Name:          "invalid full data",
		Pre:           pre,
		PostRoot:      "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
		InputMessages: []*types.SignedSSVMessage{msg},
		ExpectedError: "invalid signed message: H(data) != root",
	}
}
