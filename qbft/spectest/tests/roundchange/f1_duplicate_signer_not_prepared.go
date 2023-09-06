package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// F1DuplicateSignerNotPrepared tests not accepting f+1 speed for duplicate signer (not prev prepared)
func F1DuplicateSignerNotPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 5),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change f+1 not duplicate prepared",
		Pre:            pre,
		PostRoot:       "15d349c8d55a140c5a64d256b8fe53c34a2e3545d75f7e787fefbc03c2630291",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}
