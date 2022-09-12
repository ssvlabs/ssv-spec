package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidRoundChangeData tests a round change msg data for which Validate() != nil
func InvalidRoundChangeData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	signQBFTMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         qbft.FirstRound,
		Input:         nil,
		PreparedRound: qbft.FirstRound,
	}).Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: signQBFTMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "invalid round change data",
		Pre:              pre,
		PostRoot:         "56cee2fd474513bc56851dfbb027366f6fc3f90fe8fec4081e993b69f84e2228",
		InputMessagesSIP: msgs,
		OutputMessages:   []*qbft.SignedMessage{},
		ExpectedError:    "round change msg invalid: round change justification invalid",
	}
}
