package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigners tests future msg with multiple signers
func MultiSigners() *ControllerSyncSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	return &ControllerSyncSpecTest{
		Name: "future msgs multiple signers",
		InputMessages: []*qbft.SignedMessage{
			testingutils.MultiSignQBFTMsg(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				&qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     2,
					Round:      qbft.FirstRound,
					Identifier: identifier[:],
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				}),
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
		ExpectedError:        "invalid future msg: allows 1 signer",
	}
}
