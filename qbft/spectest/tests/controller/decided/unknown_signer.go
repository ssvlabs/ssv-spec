package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// UnknownSigner tests a decided msg with an unknown signer
func UnknownSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
			[]types.OperatorID{1, 2, 5},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     2,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "decided unknown signer",
		Pre:            pre,
		PostRoot:       "3e721f04a2a64737ec96192d59e90dfdc93f166ec9a21b88cc33ee0c43f2b26a",
		InputMessages:  msgs,
		ExpectedError:  "invalid decided msg: invalid decided msg: commit Height is wrong",
		OutputMessages: []*qbft.SignedMessage{},
	}
}
