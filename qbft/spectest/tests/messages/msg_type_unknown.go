package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
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
		Identifier: testingutils.TestingIdentifier,
		Root:       testingutils.TestingQBFTRootData,
	})

	test := tests.NewMsgSpecTest(
		"msg type unknown",
		testdoc.MessagesMsgTypeUnknownDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		"message type is invalid",
	)

	test.SetPrivateKeys(ks)

	return test
}
