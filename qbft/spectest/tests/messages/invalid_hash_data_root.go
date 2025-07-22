package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidHashDataRoot tests an invalid hash data root
func InvalidHashDataRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: testingutils.TestingIdentifier,
		Root:       testingutils.DifferentRoot,
	})

	msg.FullData = testingutils.TestingQBFTFullData

	return tests.NewMsgSpecTest(
		"invalid hash data root",
		testdoc.MessagesInvalidHashDataRootDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		"",
	)
}
