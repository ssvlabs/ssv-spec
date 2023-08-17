package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// WrongMsgType tests a non commit msg with 2f+1 signers
func WrongMsgType() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.ControllerSpecTest{
		Name: "decide wrong msg type",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingMultiSignerProposalMessageWithHeight(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						qbft.FirstHeight,
					),
				},
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
			},
		},
		ExpectedError: "could not process msg: invalid signed message: msg allows 1 signer",
	}
}
