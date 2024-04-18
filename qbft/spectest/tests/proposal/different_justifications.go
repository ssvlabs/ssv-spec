package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DifferentJustifications tests a proposal for > 1 round, prepared previously with rc justification prepares at different heights (tests the highest prepared calculation)
func DifferentJustifications() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 3
	ks4 := testingutils.Testing4SharesSet()
	ks10 := testingutils.Testing10SharesSet()

	prepareMsgs1 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks4.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks4.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks4.Shares[3], types.OperatorID(3)),
	}
	prepareMsgs2 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks4.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks4.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks4.Shares[3], types.OperatorID(3), 2),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithParams(
			ks4.Shares[1], types.OperatorID(1), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData, qbft.FirstRound,
			testingutils.MarshalJustifications(prepareMsgs1),
		),
		testingutils.TestingRoundChangeMessageWithParams(
			ks4.Shares[2], types.OperatorID(2), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData, 2,
			testingutils.MarshalJustifications(prepareMsgs2),
		),
		testingutils.TestingRoundChangeMessageWithParams(
			ks4.Shares[3], types.OperatorID(3), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData, qbft.NoRound,
			nil,
		),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks4.Shares[1], types.OperatorID(1), 3, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs2),
		),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "different proposal round change justification",
		Pre:           pre,
		PostRoot:      "eeee6a5783cec81c67ffe6a185f9afb4c775b80a815f22dbb32489f8ffa6adce",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithRound(ks10.Shares[1], types.OperatorID(1), 3),
		},
	}
}
