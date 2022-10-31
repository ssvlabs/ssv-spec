package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMsgSigTooLong tests SignedMessage len(signature) > 96
func SignedMsgSigTooLong() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}},
	})

	msg.Signature = make([]byte, 97)
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "signature too long",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusCommitMsgType),
				Data: b,
			},
		},
		ExpectedError: "message signature is invalid",
	}
}
