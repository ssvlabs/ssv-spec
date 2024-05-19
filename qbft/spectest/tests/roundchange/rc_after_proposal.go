package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// AfterProposal tests a round change msg after prev quorum + received proposal
func AfterProposal() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}

	msgs := append(rcMsgs, []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),

		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[4], types.OperatorID(4), 2),
	}...)

	return &tests.MsgProcessingSpecTest{
		Name:          "round change after proposal accepted",
		Pre:           pre,
		PostRoot:      "47c00aa37311b5dedadac5dd302e9b23635a7bbbe1166f221ee2e8773df49e4d",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
				testingutils.MarshalJustifications(rcMsgs)),
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		},
	}
}
