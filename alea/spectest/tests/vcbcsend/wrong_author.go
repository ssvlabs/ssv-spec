package vcbcsend

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongAuthor tests the receipt of a VCBCSend with wrong author
func WrongAuthor() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.VCBCSendMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, alea.FirstPriority, types.OperatorID(3)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcsend wrong author of a vcbcsend",
		Pre:           pre,
		PostRoot:      "f831e111116961c37e9c383d1e6e3532e2ec3d1513dfe995f002f7be80e64a8c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: author of VCBCSend differs from sender of the message",
		DontRunAC:     true,
	}
}
