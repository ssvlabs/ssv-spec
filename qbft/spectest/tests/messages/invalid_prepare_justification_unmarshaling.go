package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidPrepareJustificationsUnmarshalling tests unmarshalling invalid prepare justifications (during message.validate())
func InvalidPrepareJustificationsUnmarshalling() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:              qbft.ProposalMsgType,
		Height:               qbft.FirstHeight,
		Round:                qbft.FirstRound,
		Identifier:           testingutils.TestingIdentifier,
		Root:                 testingutils.DifferentRoot,
		PrepareJustification: [][]byte{{1}},
	})

	msg.FullData = testingutils.TestingQBFTFullData

	test := tests.NewMsgSpecTest(
		"invalid prepare justification unmarshalling",
		testdoc.MessagesInvalidPrepareJustificationUnmarshalingDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		"incorrect size",
		ks,
	)

	return test
}
