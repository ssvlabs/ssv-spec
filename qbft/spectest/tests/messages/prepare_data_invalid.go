package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// PrepareDataInvalid tests prepare data len == 0
func PrepareDataInvalid() *tests.MsgSpecTest {
	d := &qbft.PrepareData{}
	byts, _ := d.Encode()
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     qbft.FirstHeight,
		Round:      10,
		Identifier: []byte{1, 2, 3, 4},
		Data:       byts,
	})

	return &tests.MsgSpecTest{
		Name: "prepare data invalid",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "PrepareData data is invalid",
	}
}
