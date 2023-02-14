package vcbcrequest

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Receive tests the receipt of a VCBCRequest message
func Receive() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.VCBCSendMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, alea.FirstPriority, types.OperatorID(2)),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.VCBCFinalMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCFinalDataBytes(tests.Hash, alea.FirstPriority, tests.AggregatedMsgBytes(types.OperatorID(2), alea.FirstPriority), types.OperatorID(2)),
		}),
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
			MsgType:    alea.VCBCRequestMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCRequestDataBytes(alea.FirstPriority, types.OperatorID(2)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcrequest receive",
		Pre:           pre,
		PostRoot:      "c4b71dfcf939585b980ae36a3583aa46e4427a246dc40182cf430f5ae22df48f",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCReadyMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCReadyDataBytes(tests.Hash, alea.FirstPriority, types.OperatorID(2)),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCAnswerMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCAnswerDataBytes(tests.ProposalDataList, alea.FirstPriority, tests.AggregatedMsgBytes(types.OperatorID(2), alea.FirstPriority), types.OperatorID(2)),
			}),
		},
		DontRunAC: true,
	}
}
