package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// SecondMsg tests processing a decided msg after already receiving a decided msg
func SecondMsg() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	msgs := []*qbft.SignedMessage{
		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[4]},
			[]types.OperatorID{1, 2, 4},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "decided second msg",
		Pre:            pre,
		PostRoot:       "ed99ab91cac917c5bf9ff90eee30f21fe47d2e272d1f35d005dbdffef426ac02",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
	}
}
