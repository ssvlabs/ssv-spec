package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidValueCheck tests a proposal that doesn't pass value check
func InvalidValueCheck() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	proposeMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, &qbft.Data{
		Root:   [32]byte{1, 1, 1, 1},
		Source: testingutils.TestingInvalidValueCheck,
	}).Encode()

	return &tests.MsgProcessingSpecTest{
		Name:     "invalid proposal value check",
		Pre:      pre,
		PostRoot: "56cee2fd474513bc56851dfbb027366f6fc3f90fe8fec4081e993b69f84e2228",
		InputMessages: []*types.Message{
			{
				ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusProposeMsgType),
				Data: proposeMsgEncoded,
			},
		},
		OutputMessages: []*types.Message{},
		ExpectedError:  "proposal invalid: proposal not justified: proposal value invalid: invalid value",
	}
}
