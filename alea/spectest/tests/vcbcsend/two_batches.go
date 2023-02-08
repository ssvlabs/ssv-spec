package vcbcsend

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// TwoBatches tests the delivery of two VCBCSend msg after receiving two batch of messages
func TwoBatches() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{tests.SignedProposal1, tests.SignedProposal2, tests.SignedProposal3, tests.SignedProposal4}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcsend two batch",
		Pre:           pre,
		PostRoot:      "774edf5499f4c4573a79be4bf4b1b6243db05e7c521699110537895a0869a116",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, alea.FirstPriority, types.OperatorID(1)),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList2, alea.FirstPriority+1, types.OperatorID(1)),
			}),
		},
		DontRunAC: true,
	}
}
