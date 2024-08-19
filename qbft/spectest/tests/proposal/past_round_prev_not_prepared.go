package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PastRoundProposalPrevNotPrepared tests a valid proposal for past round (not prev prepared)
func PastRoundProposalPrevNotPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 10
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingRoundChangeMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingRoundChangeMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), qbft.FirstRound,
			testingutils.MarshalJustifications(rcMsgs)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal past round (not prev prepared)",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: past round",
	}
}
