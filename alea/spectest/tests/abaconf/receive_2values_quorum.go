package abaconf

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Receive2ValuesQuorum tests an ABAConf quorum of 2 values
func Receive2ValuesQuorum() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{}
	for opID := 2; opID <= 3; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}
	for opID := 3; opID <= 4; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}
	for opID := 2; opID <= 3; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAAuxMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAAuxDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}
	for opID := 3; opID <= 4; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAAuxMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAAuxDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}
	for opID := 2; opID <= 4; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAConfMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAConfDataBytes([]byte{0, 1}, alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "abaconf receive 2values quorum",
		Pre:           pre,
		PostRoot:      "9a811c976980492dd54caec5ebfb389c80ccf39f40a08ab0ab4c66a451f4c035",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes([]byte{1, 0}, alea.FirstRound, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound+1, alea.FirstACRound),
			}),
		},
		DontRunAC: true,
	}
}
