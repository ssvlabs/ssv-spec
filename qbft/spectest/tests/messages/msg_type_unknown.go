package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgTypeUnknown TODO<olegshmuelov> validate message type for unknown or non-exist
// MsgTypeUnknown tests Message type > 5
func MsgTypeUnknown() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: [32]byte{1, 2, 3, 4}}).Encode()

	return &tests.MsgSpecTest{
		Name: "msg type unknown",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.MsgType{0x9, 0x0, 0x0, 0x0}),
				Data: msgEncoded,
			},
		},
		ExpectedError: "message type is invalid",
	}
}
