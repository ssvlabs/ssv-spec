package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidIdentifier tests a process msg with the wrong identifier
func InvalidIdentifier() *tests.ControllerSpecTest {
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	inputData := &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}}
	signMsgEncoded, _ := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
	}, inputData).Encode()
	return &tests.ControllerSpecTest{
		Name: "invalid identifier",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: inputData,
				InputMessages: []*types.Message{
					{
						ID:   types.PopulateMsgType(identifier, types.ConsensusProposeMsgType),
						Data: signMsgEncoded,
					}},
				DecidedVal:         nil,
				ControllerPostRoot: "5a1536414abb7928a962cc82e7307b48e3d6c17da15c3f09948c20bd89d41301",
			},
		},
		ExpectedError: "invalid msg: message doesn't belong to Identifier",
	}
}
