package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRoundProposalPrevNotPrepared tests a valid proposal for past round (not prev prepared)
func PastRoundProposalPrevNotPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 10
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingRoundChangeMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingRoundChangeMessage(ks.Shares[3], types.OperatorID(3)),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), qbft.FirstRound,
			testingutils.MarshalJustifications(rcMsgs)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal past round (not prev prepared)",
		Pre:            pre,
		PostRoot:       "ed0b4ac99e52e0e2be985db854913958e62d52a4424bb77fa69fc606a9060bbd",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: past round",
	}
}
