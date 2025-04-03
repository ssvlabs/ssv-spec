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
		ExpectedRoot: "39e9ac7d177bf75fffa99354d498ffea607dc473a68b3d12b6425b15efda9c9e",
	}
}
