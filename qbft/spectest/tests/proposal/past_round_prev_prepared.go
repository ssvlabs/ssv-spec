package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRoundProposalPrevPrepared tests a valid proposal for past round (prev prepared)
func PastRoundProposalPrevPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 10

	ks := testingutils.Testing4SharesSet()
	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 6),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 6),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 6),
	}

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 8),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 8),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 8),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 8, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs)),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal past round (not prev prepared)",
		Pre:            pre,
		PostRoot:       "9be139118d79649f5843fa7bb161a665e928c7e85eb4b08b624b2a40dc8ef9a1",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: past round",
	}
}
