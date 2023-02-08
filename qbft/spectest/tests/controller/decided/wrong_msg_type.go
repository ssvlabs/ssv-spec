package decided

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// WrongMsgType tests a non commit msg with 2f+1 signers
func WrongMsgType() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide wrong msg type",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						&qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     10,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
		},
		ExpectedError: "invalid future msg: allows 1 signer",
	}
}
