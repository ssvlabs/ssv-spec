package abafinish

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// ReceiveNoQuorum tests an ABAFinish no quorum
func ReceiveNoQuorum() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{}

	signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea(tests.ProposalData1.Data),
	})
	msgs = append(msgs, signedMsg)
	signedMsg = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea(tests.ProposalData2.Data),
	})
	msgs = append(msgs, signedMsg)

	for opID := 2; opID <= 2; opID++ {
		signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAFinishMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAFinishDataBytes(byte(1), alea.FirstACRound),
		})
		msgs = append(msgs, signedMsg)
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "abafinish receive no quorum",
		Pre:           pre,
		PostRoot:      "67a0782397a8d060957fe3011e973f7f25abb7493eb1512ed46ce73c4581bcd2",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, alea.FirstPriority, types.OperatorID(1)),
			}),
		},
		DontRunAC: false,
	}
}
