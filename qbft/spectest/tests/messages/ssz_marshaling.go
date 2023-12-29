package messages

import (
	"bytes"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SSZMarshaling tests a valid marshaling of a signed message
func SSZMarshaling() tests.SpecTest {
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

	b, err := msg.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	unmarshalledMsg := &qbft.SignedMessage{}
	unmarshalledMsg.UnmarshalSSZ(b)

	root1, err := msg.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	root2, err := unmarshalledMsg.GetRoot()
	if err != nil {
		panic(err.Error())
	}

	if !bytes.Equal(root1[:], root2[:]) {
		panic("Unmarshalled message is different.")
	}

	return &tests.MsgSpecTest{
		Name: "SSZ marshalling of signed messaged",
		Messages: []*qbft.SignedMessage{
			msg,
		},
	}
}
