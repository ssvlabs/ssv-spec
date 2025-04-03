package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GetRoot tests GetRoot on SignedMessage
func GetRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessageWithParams(
		ks.OperatorKeys[1], types.OperatorID(1), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications([]*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		}),
		testingutils.MarshalJustifications([]*types.SignedSSVMessage{
			testingutils.TestingRoundChangeMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		}))

	r, _ := msg.GetRoot()

	return &tests.MsgSpecTest{
		Name: "get root",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		ExpectedRoots: [][32]byte{
			r,
		},
	}
}
