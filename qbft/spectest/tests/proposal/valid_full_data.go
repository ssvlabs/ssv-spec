package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidFullData tests the signed proposal with a valid full data field (H(full data) == root)
func ValidFullData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithIdentifierAndFullData(
			ks.Shares[1], types.OperatorID(1), []byte{1, 2, 3, 4}, testingutils.TestingQBFTFullData,
			qbft.FirstHeight),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "proposal valid full data",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1))},
	}
}
