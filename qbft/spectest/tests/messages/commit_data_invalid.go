package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CommitDataInvalid tests commit data len == 0
func CommitDataInvalid() *tests.MsgSpecTest {
	d := &qbft.CommitData{}
	byts, _ := d.Encode()
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      10,
		Identifier: []byte{1, 2, 3, 4},
		Data:       byts,
	})

	return &tests.MsgSpecTest{
		Name: "commit data invalid",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "CommitData data is invalid",
	}
}
