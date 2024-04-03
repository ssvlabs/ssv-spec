package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidHashDataRoot tests an invalid hash data root
func InvalidHashDataRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.NetworkKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Root:       testingutils.DifferentRoot,
		FullData:   testingutils.TestingQBFTFullData,
	})

	return &tests.MsgSpecTest{
		Name: "invalid hash data root",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
	}
}
