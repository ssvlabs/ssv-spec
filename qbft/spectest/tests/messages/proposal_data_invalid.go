package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// ProposalDataInvalid tests proposal data len == 0
func ProposalDataInvalid() *tests.MsgSpecTest {
	d := &qbft.ProposalData{}
	byts, _ := d.Encode()
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      10,
		Identifier: []byte{1, 2, 3, 4},
		Data:       byts,
	})

	return &tests.MsgSpecTest{
		Name: "proposal data invalid",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "ProposalData data is invalid",
	}
}
