package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SignedMsgMultiSigners tests SignedMessage with multi signers
func SignedMsgMultiSigners() *tests.MsgSpecTest {
	baseMsgID := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msgEncoded, _ := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{
			testingutils.Testing4SharesSet().Shares[1],
			testingutils.Testing4SharesSet().Shares[2],
			testingutils.Testing4SharesSet().Shares[3],
		},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}},
		}).Encode()

	return &tests.MsgSpecTest{
		Name: "multi signers",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(baseMsgID, types.DecidedMsgType),
				Data: msgEncoded,
			},
		},
	}
}
