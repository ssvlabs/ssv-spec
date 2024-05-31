package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// F1DuplicateSignerNotPrepared tests not accepting f+1 speed for duplicate signer (not prev prepared)
func F1DuplicateSignerNotPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 5),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change f+1 not duplicate prepared",
		Pre:            pre,
		PostRoot:       "c0a3241126f85aa19a5dc89fc666930185d14a7652b5ac4dbb526239cae9b23b",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
