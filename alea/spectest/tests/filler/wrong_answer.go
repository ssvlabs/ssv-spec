package filler

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
			MsgType:    alea.FillerMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.FillerDataBytes([][]*alea.ProposalData{tests.ProposalDataList2}, []alea.Priority{alea.FirstPriority}, [][]byte{tests.AggregatedMsgBytes2(types.OperatorID(1), alea.FirstPriority)}, types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "filler wrong answer",
		Pre:           pre,
		PostRoot:      "4f22c4a4a5c89a32df7844556558ab48ddc65c491d09e2c8e4a89374edf2c0d8",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: existing (priority,author) with different proposals",
		DontRunAC:     true,
	}
}
