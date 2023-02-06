package abaaux

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Duplicate tests an ABAInit received twice
func Duplicate() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.ABAAuxMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAAuxDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.ABAAuxMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAAuxDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "abaaux duplicate",
		Pre:           pre,
		PostRoot:      "a41f3acae42d3d38fc1693b4d9219e398f7b0b51e9000df432c29a137fbc009a",
		InputMessages: msgs,
		DontRunAC:     true,
	}
}
