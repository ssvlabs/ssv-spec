package processmsg

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
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
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
		},
		ExpectedError: "could not process msg: invalid signed message: proposal not justified: change round has no quorum",
	}
}
