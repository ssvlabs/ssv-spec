package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMsgNoSigners tests SignedMessage len(signers) == 0
func SignedMsgNoSigners() *tests.MsgSpecTest {
	baseMsgID := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	msg.Signers = nil
	msgEncoded, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "no signers",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(baseMsgID, types.ConsensusCommitMsgType),
				Data: msgEncoded,
			},
		},
		ExpectedError: "message signers is empty",
	}
}
