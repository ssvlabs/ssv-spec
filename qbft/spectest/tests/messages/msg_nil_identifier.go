package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgNilIdentifier TODO<olegshmuelov> find a way to validate the identifier
// MsgNilIdentifier tests Message with Identifier == nil
func MsgNilIdentifier() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID(nil, types.BNRoleAttester)
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}},
	})
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "msg identifier nil",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
				Data: b,
			},
		},
		ExpectedError: "message identifier is invalid",
	}
}
