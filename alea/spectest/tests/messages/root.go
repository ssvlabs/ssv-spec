package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// GetRoot tests GetRoot on SignedMessage
func GetRoot() *tests.MsgSpecTest {
	msg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data: testingutils.ProposalDataBytesAlea(
			[]byte{1, 2, 3, 4},
		),
	})

	r, _ := msg.GetRoot()

	return &tests.MsgSpecTest{
		Name: "get root",
		Messages: []*alea.SignedMessage{
			msg,
		},
		ExpectedRoots: [][]byte{
			r,
		},
	}
}
