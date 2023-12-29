package messages

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PrepareJustificationsUnmarshalling tests unmarshalling prepare justifications
func PrepareJustificationsUnmarshalling() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], 1, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], 2, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], 3, 2),
	}

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}

	rcMarshalled := testingutils.MarshalJustifications(rcMsgs)

	prepareMarshalled := testingutils.MarshalJustifications(prepareMsgs)

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		rcMarshalled, prepareMarshalled)

	// Assert unmarshalling is correct
	prepareUnmarshalled, err := msg.Message.GetPrepareJustifications()
	if err != nil {
		panic(err)
	}

	// Compare messages
	for idx, prepareMsg := range prepareUnmarshalled {
		root1, err := prepareMsg.GetRoot()
		if err != nil {
			panic(err)
		}
		root2, err := prepareMsgs[idx].GetRoot()
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(root1[:], root2[:]) {
			panic("Unmarshalled message is different")
		}
	}

	return &tests.MsgSpecTest{
		Name: "prepare justification unmarshalling",
		Messages: []*qbft.SignedMessage{
			msg,
		},
	}
}
