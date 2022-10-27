package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeDataInvalidJustifications tests PreparedRound != NoRound len(RoundChangeJustification) == 0
func RoundChangeDataInvalidJustifications() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	rcMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         10,
		Input:         &qbft.Data{Root: [32]byte{}, Source: []byte{1, 2, 3, 4}},
		PreparedRound: 1,
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(identifier, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
	}

	return &tests.MsgSpecTest{
		Name:          "rc prev prepared no justifications",
		Messages:      msgs,
		ExpectedError: "round change justification invalid",
	}
}
