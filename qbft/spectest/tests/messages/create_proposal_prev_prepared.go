package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CreateProposalPreviouslyPrepared tests creating a proposal msg,previously prepared
func CreateProposalPreviouslyPrepared() *tests.CreateMsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create proposal previously prepared",
		Value:      [32]byte{1, 2, 3, 4},
		RoundChangeJustifications: []*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
			testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
		},
		PrepareJustifications: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
			testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
		},
		ExpectedRoot: "15a4700a55577a524a147d348d35aaeafbd04b478ba963988e7af3647b9dbb8c",
	}
}
