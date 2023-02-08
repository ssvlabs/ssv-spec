package vcbcsend

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Batch tests the delivery of a VCBCSend msg after receiving a batch of messages
func Batch() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{tests.SignedProposal1, tests.SignedProposal2}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcsend batch",
		Pre:           pre,
		PostRoot:      "b2a407ff98d1b828229339e758543442535cf9d4fce1e68893dba0009b845c12",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, alea.FirstPriority, types.OperatorID(1)),
			}),
		},
		DontRunAC: true,
	}
}
