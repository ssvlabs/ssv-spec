package vcbcanswer

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongData tests the case in which the message has proposals that mismatch with the aggreagtedMsgBytes's proposals
func WrongData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.VCBCAnswerMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCAnswerDataBytes(tests.ProposalDataList2, alea.FirstPriority, tests.AggregatedMsgBytes(types.OperatorID(1), alea.FirstPriority), types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcanswer wrong data",
		Pre:           pre,
		PostRoot:      "d0669999d2f4f17dd4888e9602362eb73a7c961e8090c5e5ea2e5e6d5608e9cd",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: hash of proposals given doesn't match hash in the VCBCReadyData of the aggregated message",
		DontRunAC:     true,
	}
}
