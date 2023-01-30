package fillgap

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongAuthor tests the receipt of a VCBCRequest message with operatorID author invalid
func WrongAuthor() *tests.MsgProcessingSpecTest {
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
			MsgType:    alea.FillGapMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.FillGapDataBytes(types.OperatorID(5), alea.FirstPriority),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "fillgap receive wrong author",
		Pre:           pre,
		PostRoot:      "5068c0e43e247c8e6ffe24e163cee775adb9f7020a7c592309c850c8be0061e4",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: author (OperatorID) doesn't exist in Committee",
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCReadyMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCReadyDataBytes(tests.Hash, alea.FirstPriority, types.OperatorID(2)),
			}),
		},
		DontRunAC: true,
	}
}
