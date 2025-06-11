package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CreateProposalPreviouslyPrepared tests creating a proposal msg,previously prepared
func CreateProposalPreviouslyPrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create proposal previously prepared",
		Value:      [32]byte{1, 2, 3, 4},
		RoundChangeJustifications: []*types.SignedSSVMessage{
			testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], types.OperatorID(1), 2),
			testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), 2),
			testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), 2),
		},
		PrepareJustifications: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
			testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
		},
		ExpectedRoot: "e304c9709d6ba8851663b9e01626424e1960cba29f4d72475c3c65c939232544",
	}
}
