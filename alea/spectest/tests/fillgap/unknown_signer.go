package fillgap

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// UnknownSigner tests a single proposal received with an unknown signer
func UnknownSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(5), &alea.Message{
			MsgType:    alea.FillGapMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.FillGapDataBytes(types.OperatorID(1), alea.FirstPriority),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "fillgap unknown proposal signer",
		Pre:           pre,
		PostRoot:      "d0669999d2f4f17dd4888e9602362eb73a7c961e8090c5e5ea2e5e6d5608e9cd",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg signature invalid: unknown signer",
		DontRunAC:     true,
	}
}
