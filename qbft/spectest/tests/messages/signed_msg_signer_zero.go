package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SignedMessageSigner0 tests SignedMessage signer == 0
func SignedMessageSigner0() *tests.MsgSpecTest {
	baseMsgID := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	msgEncoded, _ := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 0},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}},
		}).Encode()

	return &tests.MsgSpecTest{
		Name: "signer 0",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(baseMsgID, types.DecidedMsgType),
				Data: msgEncoded,
			},
		},
		ExpectedError: "signer ID 0 not allowed",
	}
}
