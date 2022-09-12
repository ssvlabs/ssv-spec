package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PastRound tests a commit msg with past round, should process but not decide
func PastRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	signMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  5,
		Input:  []byte{1, 2, 3, 4},
	})
	pre.State.ProposalAcceptedForCurrentRound = signMsg
	pre.State.Round = 5
	signMsg.Message.Round = 2
	signMsgEncoded, _ := signMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "commit past round",
		Pre:              pre,
		PostRoot:         "500018acb258bca93ea5fd6fc2c5f649b07f9a6f6add09fb626af605b780a057",
		InputMessagesSIP: msgs,
		ExpectedError:    "commit msg invalid: commit round is wrong",
	}
}
