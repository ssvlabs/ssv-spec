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
	commitMsgEncoded, _ := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[3]},
		[]types.OperatorID{1, 2, 3},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  []byte{1, 2, 3, 4},
		}).Encode()
	commitMsgEncoded2, _ := testingutils.MultiSignQBFTMsg(
		[]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2], testingutils.Testing4SharesSet().Shares[4]},
		[]types.OperatorID{1, 2, 4},
		&qbft.Message{
			Height: qbft.FirstHeight,
			Round:  qbft.FirstRound,
			Input:  []byte{1, 2, 3, 4},
		}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: commitMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: commitMsgEncoded2,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "decided second msg",
		Pre:            pre,
		PostRoot:       "ed99ab91cac917c5bf9ff90eee30f21fe47d2e272d1f35d005dbdffef426ac02",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
	}
}
