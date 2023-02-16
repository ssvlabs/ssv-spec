package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CurrentRoundPrevPrepared tests a > first round proposal prev prepared
func CurrentRoundPrevPrepared() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 10

	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 8),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 8),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 8),
	}

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 10),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 10),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(
			ks.Shares[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "proposal happy flow round > 1 (prev prepared)",
		Pre:           pre,
		PostRoot:      "06cae71db44b25f0a53e190ad833186339629de2e65599a9c3dab63127277c48",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], 1, 10),
		},
	}
}
