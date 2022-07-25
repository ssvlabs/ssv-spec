package runnertests

import (
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	index := 4
	//ks := testingutils.Testing4SharesSet()
	dataSet := testutils.TestDepositSignDataSetFourOperators()
	//pre := testutils.TestRunner(dataSet, *ks, uint64(index))
	//node := testutils.TestNode(dataSet, uint64(index))
	messages := map[types.OperatorID]*dkgtypes.Message{
		1: testutils.SignDKGMsg(dataSet.DKGOperators[1].SK, dataSet.ParsedSignedDepositDataMessage(1)).(*dkgtypes.Message),
		2: testutils.SignDKGMsg(dataSet.DKGOperators[2].SK, dataSet.ParsedSignedDepositDataMessage(2)).(*dkgtypes.Message),
		3: testutils.SignDKGMsg(dataSet.DKGOperators[3].SK, dataSet.ParsedSignedDepositDataMessage(3)).(*dkgtypes.Message),
	}
	prasedMessages := func(msgs map[types.OperatorID]*dkgtypes.Message) map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage {
		out := map[types.OperatorID]*dkgtypes.ParsedSignedDepositDataMessage{}
		for _, m := range msgs {
			parsed := &dkgtypes.ParsedSignedDepositDataMessage{}
			parsed.FromBase(m)
			out[types.OperatorID(m.Header.Sender)] = parsed
		}
		return out
	}(messages)

	return &MsgProcessingSpecTest{
		Name:     "happy flow",
		Operator: dataSet.Operator(types.OperatorID(index)),
		Config:   testutils.TestConfig(dataSet, uint64(index)),
		KeySet:   &dataSet.TestKeySet,
		Output:   prasedMessages,
		Messages: []*dkgtypes.Message{
			testutils.SignDKGMsg(dataSet.DKGOperators[1].SK, dataSet.ParsedInitMessage(1)).(*dkgtypes.Message),
			testutils.SignDKGMsg(dataSet.DKGOperators[1].SK, testutils.PlaceholderMessage()).(*dkgtypes.Message),
			testutils.SignDKGMsg(dataSet.DKGOperators[2].SK, dataSet.ParsedPartialSigMessage(2)).(*dkgtypes.Message),
			testutils.SignDKGMsg(dataSet.DKGOperators[3].SK, dataSet.ParsedPartialSigMessage(3)).(*dkgtypes.Message),
			messages[1],
			messages[2],
			messages[3],
		},
		Outgoing: []*dkgtypes.Message{
			dataSet.ParsedPartialSigMessage(4),
			dataSet.ParsedSignedDepositDataMessage(4),
		},
	}
}
