package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// Invalid tests decided msg where msg.validate() != nil
func Invalid() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     10,
			Round:      qbft.FirstRound,
			Identifier: identifier[:],
			Data:       []byte{1, 2, 3, 4},
		})
	msg.Signers = []types.OperatorID{}
	return &tests.ControllerSpecTest{
		Name: "decide invalid msg",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					msg,
				},
				ControllerPostRoot: "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
			},
		},
		ExpectedError: "invalid future msg: invalid decided msg: message signers is empty",
	}
}
