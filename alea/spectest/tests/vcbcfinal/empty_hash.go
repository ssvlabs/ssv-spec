package vcbcfinal

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// EmptyHash tests a proposal msg received with the wrong height
func EmptyHash() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.VCBCFinalMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCFinalDataBytes([]byte{}, alea.FirstPriority, tests.AggregatedMsgBytes(types.OperatorID(1), alea.FirstPriority), types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcfinal empty hash",
		Pre:           pre,
		PostRoot:      "f831e111116961c37e9c383d1e6e3532e2ec3d1513dfe995f002f7be80e64a8c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: VCBCFinalData invalid: VCBCFinalData: empty hash",
		DontRunAC:     true,
	}
}
