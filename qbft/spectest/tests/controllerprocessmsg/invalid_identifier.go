package controllerprocessmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidIdentifier tests a process msg with the wrong identifier
func InvalidIdentifier() *tests.ControllerSpecTest {
	return &tests.ControllerSpecTest{
		Name: "invalid identifier",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], 1, &qbft.Message{
						MsgType:    qbft.ProposalMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: []byte{1, 2, 3, 4},
						Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
					}),
				},
				DecidedVal:         nil,
				ControllerPostRoot: "f28cfa54f7993a21ebe3a46fde514e387ad4ff9a3be197b656c4c6e8dbb124de",
			},
		},
		ExpectedError: "invalid msg: message doesn't belong to Identifier",
	}
}
