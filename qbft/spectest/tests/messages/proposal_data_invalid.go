package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ProposalDataInvalid tests proposal data len == 0
func ProposalDataInvalid() *tests.MsgSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  10,
		Input:  &qbft.Data{},
	}).Encode()

	return &tests.MsgSpecTest{
		Name: "proposal data invalid",
		Messages: []*types.Message{
			{
				ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
				Data: msgEncoded,
			},
		},
		ExpectedError: "message input data is invalid",
	}
}
