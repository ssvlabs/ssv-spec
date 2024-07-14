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

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
	}

	msgs := append(rcMsgs, []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),

		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[4], types.OperatorID(4), 2),
	}...)

	return &tests.MsgProcessingSpecTest{
		Name:          "round change after proposal accepted",
		Pre:           pre,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingProposalMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2,
				testingutils.MarshalJustifications(rcMsgs)),
			testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
		},
	}
}
