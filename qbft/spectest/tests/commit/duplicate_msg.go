package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// DuplicateMsg tests a duplicate commit msg processing
func DuplicateMsg() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	baseMsgId := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	signMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
	})
	pre.State.ProposalAcceptedForCurrentRound = signMsg
	signMsgEncoded, _ := signMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
		{
			ID:   types.PopulateMsgType(baseMsgId, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "duplicate commit message",
		Pre:           pre,
		PostRoot:      "952dc40814611c59abfa2cbc0445c30cc54646da48ccdac01d8b48943770c569",
		InputMessages: msgs,
	}
}
