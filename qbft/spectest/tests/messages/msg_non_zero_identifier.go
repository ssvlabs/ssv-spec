package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgNonZeroIdentifier tests Message with len(Identifier) == 0
func MsgNonZeroIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{},
		Root:       testingutils.TestingQBFTRootData,
	})

	return &tests.MsgSpecTest{
		Name: "msg identifier len == 0",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "message identifier is invalid",
	}
}
