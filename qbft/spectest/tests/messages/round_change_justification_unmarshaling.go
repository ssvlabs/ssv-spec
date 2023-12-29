package messages

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangeJustificationsUnmarshalling tests unmarshalling round change justifications
func RoundChangeJustificationsUnmarshalling() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], 1, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], 2, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], 3, 2),
	}

	rcMarshalled := testingutils.MarshalJustifications(rcMsgs)

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		rcMarshalled, nil)

	// Assert unmarshalling is correct
	rcUnmarshalled, err := msg.Message.GetRoundChangeJustifications()
	if err != nil {
		panic(err)
	}

	// Compare messages
	for idx, rcMsg := range rcUnmarshalled {
		root1, err := rcMsg.GetRoot()
		if err != nil {
			panic(err)
		}
		root2, err := rcMsgs[idx].GetRoot()
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(root1[:], root2[:]) {
			panic("Unmarshalled message is different")
		}
	}

	return &tests.MsgSpecTest{
		Name: "round change justification unmarshalling",
		Messages: []*qbft.SignedMessage{
			msg,
		},
	}
}
