package abafinish

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Receive tests an ABAFinish msg
func Receive() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.ABAFinishMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAFinishDataBytes(byte(1), alea.FirstACRound),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "abafinish receive",
		Pre:           pre,
		PostRoot:      "9c8df16a26dff87d7408be03573530d843816ca5b2909643b943ba57abaa088e",
		InputMessages: msgs,
		DontRunAC:     true,
	}
}
