package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MarshalJustificationsWithFullData tests marshalling justifications with full data (should omit it)
func MarshalJustificationsWithFullData() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	encodedRCMsg, _ := testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], 1, 2).MarshalSSZ()
	encodedPrepareMsg, _ := testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)).MarshalSSZ()

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		[][]byte{encodedRCMsg}, [][]byte{encodedPrepareMsg})

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "marshal justifications with full data",
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
