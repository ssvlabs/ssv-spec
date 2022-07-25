package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	dkgtypes "github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/types"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	index := 1
	//ks := testingutils.Testing4SharesSet()
	dataSet := testutils.TestDepositSignDataSetFourOperators()
	//pre := testutils.TestRunner(dataSet, *ks, uint64(index))
	//node := testutils.TestNode(dataSet, uint64(index))
	messages := map[types.OperatorID]*dkgtypes.Message{
		2: testutils.SignDKGMsg(dataSet.DKGOperators[2].SK, dataSet.ParsedSignedDepositDataMessage(2)).(*dkgtypes.Message),
		3: testutils.SignDKGMsg(dataSet.DKGOperators[3].SK, dataSet.ParsedSignedDepositDataMessage(3)).(*dkgtypes.Message),
		4: testutils.SignDKGMsg(dataSet.DKGOperators[4].SK, dataSet.ParsedSignedDepositDataMessage(4)).(*dkgtypes.Message),
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
		Name:          "happy flow",
		Operator:      dataSet.Operator(types.OperatorID(index)),
		LocalKeyShare: dataSet.MakeLocalKeyShare(1),
		KeySet:        &dataSet.TestKeySet,
		Output:        prasedMessages,
		Messages: []*dkgtypes.Message{
			testutils.SignDKGMsg(dataSet.DKGOperators[1].SK, dataSet.ParsedInitMessage(1)).(*dkgtypes.Message),
			testutils.SignDKGMsg(dataSet.DKGOperators[1].SK, testutils.PlaceholderMessage()).(*dkgtypes.Message),
			testutils.SignDKGMsg(dataSet.DKGOperators[2].SK, dataSet.ParsedPartialSigMessage(2)).(*dkgtypes.Message),
			testutils.SignDKGMsg(dataSet.DKGOperators[3].SK, dataSet.ParsedPartialSigMessage(3)).(*dkgtypes.Message),
			messages[2],
			messages[3],
			messages[4],
		},
		Outgoing: []*dkgtypes.Message{
			dataSet.ParsedPartialSigMessage(types.OperatorID(index)),
			dataSet.ParsedSignedDepositDataMessage(types.OperatorID(index)),
		},
	}
}
