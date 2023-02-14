package abafinish

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// InvalidVote tests an ABAFinish msg with an invalid vote
func InvalidVote() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.ABAFinishMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAFinishDataBytes(byte(2), alea.FirstACRound),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "abafinish invalid vote",
		Pre:           pre,
		PostRoot:      "d0669999d2f4f17dd4888e9602362eb73a7c961e8090c5e5ea2e5e6d5608e9cd",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: ABAFinishData invalid: ABAFinishData: vote not 0 or 1",
		DontRunAC:     true,
	}
}
