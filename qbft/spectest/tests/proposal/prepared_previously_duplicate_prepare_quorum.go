package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreparedPreviouslyDuplicatePrepareQuorum tests a proposal for > 1 round, prepared previously with quorum of prepared msgs (2 of which are duplicates)
func PreparedPreviouslyDuplicatePrepareQuorum() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[3], types.OperatorID(3), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate prepare msg justification quorum",
		Pre:           pre,
		PostRoot:      "a07a5d1dc880a573dfc1e45ba61ac636d275138060c2742157b8981148da62e1",
		InputMessages: append(msgs, testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2)),
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		},
	}
}
