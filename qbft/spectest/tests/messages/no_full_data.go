package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HasFullData tests signed message with full data
func HasFullData() *tests.MultiMsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	signMsg := func(msgType qbft.MessageType, fullData []byte) *qbft.SignedMessage {
		ret := testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    msgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Root:       testingutils.TestingQBFTRootData,
		})
		ret.FullData = fullData
		return ret
	}

	return &tests.MultiMsgSpecTest{
		Name: "has full data",
		Tests: []*tests.MsgSpecTest{
			{
				Name: "proposal",
				Messages: []*qbft.SignedMessage{
					signMsg(qbft.ProposalMsgType, testingutils.TestingQBFTFullData),
				},
			},
			{
				Name: "round change",
				Messages: []*qbft.SignedMessage{
					signMsg(qbft.RoundChangeMsgType, testingutils.TestingQBFTFullData),
				},
			},
			{
				Name: "prepare",
				Messages: []*qbft.SignedMessage{
					signMsg(qbft.PrepareMsgType, testingutils.TestingQBFTFullData),
				},
				ExpectedError: "full data should be nil",
			},
			{
				Name: "commit",
				Messages: []*qbft.SignedMessage{
					signMsg(qbft.CommitMsgType, testingutils.TestingQBFTFullData),
				},
				ExpectedError: "full data should be nil",
			},
		},
	}
}
