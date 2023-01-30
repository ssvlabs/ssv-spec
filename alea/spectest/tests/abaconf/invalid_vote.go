package abaconf

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// InvalidVote tests an ABAConf with invalid vote
func InvalidVote() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.ABAConfMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAConfDataBytes([]byte{0, 2}, alea.FirstRound, alea.FirstACRound),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "abaconf invalid vote",
		Pre:           pre,
		PostRoot:      "84ecec5237cd4c1ca3ce3044e04e792a8abed2d470f29e1dd9416ac00511eec2",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: ABAConfData invalid: ABAConfData: vote not 0 or 1",
		DontRunAC:     true,
	}
}
