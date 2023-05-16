package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CreateProposalNotPreviouslyPrepared tests creating a proposal msg, non-first round and not previously prepared
func CreateProposalNotPreviouslyPrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create proposal not previously prepared",
		Value:      []byte{1, 2, 3, 4},
		RoundChangeJustifications: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
			testingutils.TestingProposalMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
			testingutils.TestingProposalMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
		},
		ExpectedSSZRoot: "4a62e3fea5a175243591ebc64728be36ef395c3cd6b2f31e23afe5c7ce1b9811",
	}
}
