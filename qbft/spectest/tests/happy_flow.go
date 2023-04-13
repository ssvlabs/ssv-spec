package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	
	return &MsgProcessingSpecTest{
		Name:     "happy flow",
		Pre:      pre,
		PostRoot: happyFlowStateComparison().Register().Root(),
		InputMessages: []*qbft.SignedMessage{
			testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),

			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
			testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),

			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[2], types.OperatorID(2)),
			testingutils.TestingCommitMessage(ks.Shares[3], types.OperatorID(3)),
		},
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
		},
	}
}
