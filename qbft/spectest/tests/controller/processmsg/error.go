package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgError tests a process msg returning an error
func MsgError() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "process msg error",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], 1, &qbft.Message{
						MsgType:    qbft.ProposalMsgType,
						Height:     qbft.FirstHeight,
						Round:      100,
						Identifier: identifier[:],
						Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
					}),
				},
				ControllerPostRoot: "3cf8cd7d050943781b25a611cc7b51bc051608f316556a843936b62c276984cb",
			},
		},
		ExpectedError: "could not process msg: proposal invalid: proposal not justified: change round has no quorum",
	}
}
