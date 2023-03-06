package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MarshalJustifications tests marshalling justifications
func MarshalJustifications() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], 1, 2),
	}

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
	}

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs))

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "marshal justifications",
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
