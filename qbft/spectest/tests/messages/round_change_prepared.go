package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// RoundChangePrepared tests a round change prepared return value
func RoundChangePrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}

	prepareMarshalled := testingutils.MarshalJustifications(prepareMsgs)

	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData, 1, prepareMarshalled)

	return &tests.MsgSpecTest{
		Name: "round change prepared",
		Messages: []*qbft.SignedMessage{
			msg,
		},
	}
}
