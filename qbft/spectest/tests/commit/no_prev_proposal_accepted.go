package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoPrevAcceptedProposal tests a commit msg received without a previous accepted proposal
func NoPrevAcceptedProposal() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = nil
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "no previous accepted proposal",
		Pre:           pre,
		PostRoot:      "763f8d4e449334a14d6c2c3c1ba7458a210b963fa54ce95fed991da5e9847ebd",
		InputMessages: msgs,
	}
}
