package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgDataNonZero tests len(data) == 0
func MsgDataNonZero() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{},
	})

	e, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "msg data len 0",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
				Data: e,
			},
		},
		ExpectedError: "message input data is invalid",
	}
}
