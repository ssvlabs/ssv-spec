package queue

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/queue"
)

func CurrentHeightFirst() *MessagePriorityPairTest {
	return &MessagePriorityPairTest{
		Name: "Current height before higher height",
		State: &queue.State{
			HasRunningInstance: false,
			Height:             100,
			Slot:               64,
			Quorum:             4,
		},
		A: mockMessage{Consensus: &mockConsensusMessage{Height: 100, Type: qbft.PrepareMsgType}},
		B: mockMessage{Consensus: &mockConsensusMessage{Height: 101, Type: qbft.PrepareMsgType}},
	}
}
