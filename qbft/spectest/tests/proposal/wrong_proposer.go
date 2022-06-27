package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongProposer tests a proposal by the wrong proposer
func WrongProposer() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "wrong proposer",
		Pre:           pre,
		PostRoot:      "3e721f04a2a64737ec96192d59e90dfdc93f166ec9a21b88cc33ee0c43f2b26a",
		InputMessages: msgs,
		ExpectedError: "proposal invalid: proposal leader invalid",
	}
}
