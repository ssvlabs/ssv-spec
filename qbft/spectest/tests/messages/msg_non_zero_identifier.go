package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgNonZeroIdentifier tests Message with len(Identifier) == 0
func MsgNonZeroIdentifier() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{}, types.BNRoleAttester)
	msgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{Root: [32]byte{1, 2, 3, 4}}).Encode()

	return &tests.MsgSpecTest{
		Name: "msg identifier len == 0",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
				Data: msgEncoded,
			},
		},
		ExpectedError: "message identifier is invalid",
	}
}
