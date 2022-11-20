package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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
				DecidedVal:         nil,
				ControllerPostRoot: "3cf8cd7d050943781b25a611cc7b51bc051608f316556a843936b62c276984cb",
			},
		},
		ExpectedError: "invalid msg: message doesn't belong to Identifier",
	}
}
