package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgTypeUnknown tests Message type > 5
func MsgTypeUnknown() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    4,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Root:       testingutils.TestingQBFTRootData,
	})

	return &tests.MsgSpecTest{
		Name: "msg type unknown",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "message type is invalid",
	}
}
