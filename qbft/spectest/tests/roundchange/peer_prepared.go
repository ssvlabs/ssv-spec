package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PeerPrepared tests a round change quorum where a peer is the only one prepared
func PeerPrepared() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change peer prepared",
		Pre:           pre,
		PostRoot:      "7d6e2db3c1a0a9f0536c40b2e9f2d7d41cb89f102209f3f080bf3e44c89a786f",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
				testingutils.TestingQBFTRootData,
				testingutils.MarshalJustifications(msgs), testingutils.MarshalJustifications(prepareMsgs)),
		},
	}
}
