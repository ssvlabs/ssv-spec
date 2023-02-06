package abaconf

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// ReceiveNoQuorum tests an ABAConf no quorum
func ReceiveNoQuorum() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{}
	for opID := 2; opID <= 4; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}
	for opID := 2; opID <= 4; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAAuxMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAAuxDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}
	for opID := 2; opID <= 2; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAConfMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAConfDataBytes([]byte{1}, alea.FirstRound, alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "abaconf receive no quorum",
		Pre:           pre,
		PostRoot:      "086e603d09c40631f588acc33f48b7036b65101972d0d47290ef303b557a9d9f",
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
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes([]byte{1}, alea.FirstRound, alea.FirstACRound),
			}),
		},
		DontRunAC: true,
	}
}
