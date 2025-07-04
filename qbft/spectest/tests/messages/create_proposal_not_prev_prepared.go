package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateProposalNotPreviouslyPrepared tests creating a proposal msg, non-first round and not previously prepared
func CreateProposalNotPreviouslyPrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create proposal not previously prepared",
		Value:      [32]byte{1, 2, 3, 4},
		Round:      10,
		RoundChangeJustifications: []*types.SignedSSVMessage{
			testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
			testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
			testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		},
		ExpectedRoot: "6a2917ae827e875a646e88ebb1d483a0a99e4f321e7f063138e99a7e7b08794e",
	}
}
