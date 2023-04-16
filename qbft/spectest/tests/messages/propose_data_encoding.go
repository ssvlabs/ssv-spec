package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ProposeDataEncoding tests encoding ProposalData
func ProposeDataEncoding() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications([]*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		}),
		testingutils.MarshalJustifications([]*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessage(ks.Shares[1], types.OperatorID(1)),
		}))

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "propose data encoding",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][32]byte{
			r,
		},
	}
}
