package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// MsgNonZeroIdentifier tests Message with len(Identifier) == 0
func MsgNonZeroIdentifier() *tests.MsgSpecTest {
	msg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{},
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})

	return &tests.MsgSpecTest{
		Name: "msg identifier len == 0",
		Messages: []*alea.SignedMessage{
			msg,
		},
		ExpectedError: "message identifier is invalid",
	}
}
