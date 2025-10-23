package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgNilIdentifier tests Message with Identifier == nil
func MsgNilIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: nil,
		Root:       testingutils.TestingQBFTRootData,
	})

	test := tests.NewMsgSpecTest(
		"msg identifier nil",
		testdoc.MessagesMsgNilIdentifierDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		types.MessageIdentifierInvalidErrorCode,
		ks,
	)

	return test
}
