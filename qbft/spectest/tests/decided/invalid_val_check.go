package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// InvalidValCheckData tests a decided message with invalid decided data
func InvalidValCheckData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{
				testingutils.Testing4SharesSet().Shares[1],
				testingutils.Testing4SharesSet().Shares[2],
				testingutils.Testing4SharesSet().Shares[3],
			},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.CommitDataBytes(testingutils.TestingInvalidValueCheck),
			}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "decided invalid val check data",
		Pre:            pre,
		PostRoot:       "3e721f04a2a64737ec96192d59e90dfdc93f166ec9a21b88cc33ee0c43f2b26a",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid decided msg: decided value invalid: invalid value",
	}
}
