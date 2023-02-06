package abainit

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// ReceiveTwoQuorum tests two ABAInit quorum receipt
func ReceiveTwoQuorum() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound+1, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound+1, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound+1, alea.FirstACRound),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound+1),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound+1),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound+1),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "abainit receive two quorum",
		Pre:           pre,
		PostRoot:      "30f9691016fcff6212473e1b1e37a9022a9055aae4c4b9db1e9f459a7e4d28be",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
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
				Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound+1, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(byte(0), alea.FirstRound+1, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound+1),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(byte(0), alea.FirstRound, alea.FirstACRound+1),
			}),
		},
		DontRunAC: true,
	}
}
