package abaconf

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Receive tests an ABAConf msg
func Receive() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.ABAConfMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAConfDataBytes([]byte{0, 1}, alea.FirstRound, alea.FirstACRound),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "abaconf receive",
		Pre:           pre,
		PostRoot:      "ef09cca4f080c2ce7b1994c52643f9d2d7b126e0b20576677ab044d7ee6e89fd",
		InputMessages: msgs,
		DontRunAC:     true,
	}
}
