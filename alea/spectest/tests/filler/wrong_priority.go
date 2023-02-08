package filler

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongPriority tests the case in which the message has a priority that mismatch with the aggreagtedMsgBytes's priority
func WrongPriority() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.FillerMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.FillerDataBytes([][]*alea.ProposalData{tests.ProposalDataList}, []alea.Priority{alea.FirstPriority + 1}, [][]byte{tests.AggregatedMsgBytes(types.OperatorID(1), alea.FirstPriority)}, types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "filler wrong priority",
		Pre:           pre,
		PostRoot:      "f831e111116961c37e9c383d1e6e3532e2ec3d1513dfe995f002f7be80e64a8c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: priority given doesn't match priority in the VCBCReadyData of the aggregated message",
		DontRunAC:     true,
	}
}
