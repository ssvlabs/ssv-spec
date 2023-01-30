package vcbcsend

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
			MsgType:    alea.VCBCSendMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, alea.FirstPriority, types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "unknown proposal signer",
		Pre:           pre,
		PostRoot:      "84ecec5237cd4c1ca3ce3044e04e792a8abed2d470f29e1dd9416ac00511eec2",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: msg signature invalid: unknown signer",
		DontRunAC:     true,
	}
}
