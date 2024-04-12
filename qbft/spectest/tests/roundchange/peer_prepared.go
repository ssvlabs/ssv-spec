package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PeerPrepared tests a round change quorum where a peer is the only one prepared
func PeerPrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change peer prepared",
		Pre:           pre,
		PostRoot:      "d7a8a73bde54510216b5ead559cdbafaac7b6b9dd83bc324f20f258ac3b5cbc4",
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
				testingutils.TestingQBFTRootData,
				testingutils.MarshalJustifications(msgs), testingutils.MarshalJustifications(prepareMsgs)),
		},
	}
}
