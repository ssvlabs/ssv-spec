package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// UnmarshalJustifications tests unmarshalling justifications
func UnmarshalJustifications() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndFullData(ks.Shares[1], 1, qbft.FirstRound, testingutils.TestingQBFTFullData),
		testingutils.TestingRoundChangeMessage(ks.Shares[2], 2),
		testingutils.TestingRoundChangeMessage(ks.Shares[3], 3),
	}
	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),
	}

	encodedRCs, err := qbft.MarshalJustifications(rcMsgs)
	if err != nil {
		panic(err)
	}
	encodedPrepares, err := qbft.MarshalJustifications(prepareMsgs)
	if err != nil {
		panic(err)
	}

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		encodedRCs, encodedPrepares)

	r, err := msg.GetRoot()
	if err != nil {
		panic(err)
	}

	b, err := msg.Encode()
	if err != nil {
		panic(err)
	}

	return &tests.MsgSpecTest{
		Name: "unmarshal justifications",
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
