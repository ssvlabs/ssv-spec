package processmsg

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// InvalidIdentifier tests a process msg with the wrong identifier
func InvalidIdentifier() *tests.ControllerSpecTest {
	share := testingutils.Testing4SharesSet().Shares[1]
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
	}
	return &tests.ControllerSpecTest{
		Name: "invalid identifier",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.SignQBFTMsg(share, 1, msg),
				},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
				ControllerPostRoot: "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
			},
		},
		ExpectedError: "invalid msg: message doesn't belong to Identifier",
	}
}
