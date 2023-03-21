package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CreateProposalNotPreviouslyPrepared tests creating a proposal msg, non-first round and not previously prepared
func CreateProposalNotPreviouslyPrepared() *tests.CreateMsgSpecTest {
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
		ExpectedRoot: "45f15e2dc3a0c1bd13aa977ead44abfabab094864db4a77dcd0cd9a30e87c3bb",
	}
}
