package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidPrepareJustificationsUnmarshalling tests unmarshalling invalid prepare justifications (during message.validate())
func InvalidPrepareJustificationsUnmarshalling() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.NetworkKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:              qbft.ProposalMsgType,
		Height:               qbft.FirstHeight,
		Round:                qbft.FirstRound,
		Identifier:           []byte{1, 2, 3, 4},
		Root:                 testingutils.DifferentRoot,
		PrepareJustification: [][]byte{{1}},
		FullData:             testingutils.TestingQBFTFullData,
	})

	return &tests.MsgSpecTest{
		Name: "invalid prepare justification unmarshalling",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedError: "incorrect size",
	}
}
