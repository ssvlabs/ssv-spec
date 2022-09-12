package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgDataNil TODO<olegshmuelov> give a proper name
// MsgDataNil tests data == nil
func MsgDataNil() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  nil,
	})

	e, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "msg data nil",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
				Data: e,
			},
		},
		ExpectedError: "message input data is invalid",
	}
}
