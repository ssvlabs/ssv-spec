package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidHashDataRoot tests a valid hash data root
func ValidHashDataRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: testingutils.TestingIdentifier,
		Root:       testingutils.TestingQBFTRootData,
	})

	msg.FullData = testingutils.TestingQBFTFullData

	return tests.NewMsgSpecTest(
		"valid hash data root",
		testdoc.MessagesValidHashDataRootDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		"",
	)
}
