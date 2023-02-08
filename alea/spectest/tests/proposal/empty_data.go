package proposal

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// EmptyData tests a proposal msg received with the wrong height
func EmptyData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()
	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.ProposalMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ProposalDataBytes([]byte{}, nil, nil),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "empty data",
		Pre:           pre,
		PostRoot:      "f831e111116961c37e9c383d1e6e3532e2ec3d1513dfe995f002f7be80e64a8c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: proposalData invalid: ProposalData: data is invalid",
		DontRunAC:     true,
	}
}
