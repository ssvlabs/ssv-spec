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
		PostRoot:       "9661e3eea078696a8df907b06c96e646de3daf97c7d87f8391a4f010f30dbe4c",
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
	}
}
