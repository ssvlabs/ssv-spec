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
		Value:      [32]byte{1, 2, 3, 4},
		RoundChangeJustifications: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
			testingutils.TestingProposalMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
			testingutils.TestingProposalMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
		},
		ExpectedRoot: "f2f327aa71008b644cb5d01388f9136f02a00d4a00a25f522538f15e26dc4103",
	}
}
