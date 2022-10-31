package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidProposalData tests proposal data for which proposalData.validate() != nil
func InvalidProposalData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	signMsgInvalidEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  &qbft.Data{},
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
			Data: signMsgInvalidEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "invalid proposal data",
		Pre:            pre,
		PostRoot:       "56cee2fd474513bc56851dfbb027366f6fc3f90fe8fec4081e993b69f84e2228",
		InputMessages:  msgs,
		OutputMessages: []*types.Message{},
		ExpectedError:  "invalid signed message: message input data is invalid",
	}
}
