package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidCommitData tests commit data for which commitData.validate() != nil
func InvalidCommitData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})

	signMsgInvalidEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  nil,
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgInvalidEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "invalid commit data",
		Pre:           pre,
		PostRoot:      "0940da0c5926c8f8145e7e980ac09119e429a5208165e6dbe35d7a5c45c07555",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: message input data is invalid",
	}
}
