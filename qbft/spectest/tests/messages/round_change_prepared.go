package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RoundChangePrepared tests a round change prepared return value
func RoundChangePrepared() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}

	prepareMarshalled := testingutils.MarshalJustifications(prepareMsgs)

	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData, 1, prepareMarshalled)

	msgRoot, err := msg.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	encodedMsg, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}

	test := &tests.MsgSpecTest{
		Name: "round change prepared",
		Messages: []*types.SignedSSVMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			encodedMsg,
		},
		ExpectedRoots: []types.ExpectedRoot{
			msgRoot,
		},
	}

	test.SetPrivateKeys(ks)
	return test
}
