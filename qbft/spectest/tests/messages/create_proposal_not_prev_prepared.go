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
		ExpectedRoot: "89f88f24f9aa3496d925cb869b9de51b3be1379f1957fff15a0ba0ae83e241de",
	}
}
