package vcbcanswer

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongAnswer tests the case in which the answer to the priority have different values than the received ones
func WrongAnswer() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.VCBCSendMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, alea.FirstPriority, types.OperatorID(1)),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.VCBCAnswerMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCAnswerDataBytes(tests.ProposalDataList2, alea.FirstPriority, tests.AggregatedMsgBytes2(types.OperatorID(1), alea.FirstPriority), types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcanswer wrong answer",
		Pre:           pre,
		PostRoot:      "1766d8a4d592d6e2a46b7d5a9236bbd9ef8f0918f5eea88dcc5ad8fe102fbbc1",
		InputMessages: msgs,
		ExpectedError: "answer has different proposals than stores ones",
		DontRunAC:     true,
	}
}
