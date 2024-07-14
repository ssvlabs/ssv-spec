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
		RoundChangeJustifications: []*types.SignedSSVMessage{
			testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
			testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
			testingutils.TestingProposalMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		},
		ExpectedRoot: "b6e2b2b60655bed0f127e41d2af790fdb7d44a52282dba4888dcc9057db02353",
	}
}
