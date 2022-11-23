package queue

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/queue"
	"github.com/bloxapp/ssv-spec/ssv"
)

func NoRunningInstance() *MessagePrioritySliceTest {
	return &MessagePrioritySliceTest{
		Name: "No running instance",
		State: &queue.State{
			HasRunningInstance: false,
			Height:             100,
			Slot:               64,
			Quorum:             4,
		},
		Permutations: [][]int{
			{12, 4, 2, 6, 10, 0, 3, 11, 7, 5, 8, 1, 9},
			{8, 0, 2, 11, 1, 6, 5, 12, 9, 4, 10, 3, 7},
			{5, 12, 1, 8, 0, 2, 9, 11, 3, 6, 4, 7, 10},
		},
		Messages: []mockMessage{
			// 1. Current height/slot:
			// 1.1. Consensus
			// 1.1. Pre-consensus
			{NonConsensus: &mockNonConsensusMessage{Slot: 64, Type: ssv.SelectionProofPartialSig}},
			// 1.2. Post-consensus
			{NonConsensus: &mockNonConsensusMessage{Slot: 64, Type: ssv.PostConsensusPartialSig}},
			// 1.3.1. Consensus/Prepare
			{Consensus: &mockConsensusMessage{Height: 100, Type: qbft.PrepareMsgType}},
			// 1.3.2. Consensus/Proposal
			{Consensus: &mockConsensusMessage{Height: 100, Type: qbft.ProposalMsgType}},
			// 1.3.3. Consensus/Commit
			{Consensus: &mockConsensusMessage{Height: 100, Type: qbft.CommitMsgType}},
			// 1.3.4. Consensus/<Other>
			{Consensus: &mockConsensusMessage{Height: 100, Type: qbft.RoundChangeMsgType}},

			// 2. Higher height/slot:
			// 2.1 Decided
			{Consensus: &mockConsensusMessage{Height: 101, Decided: true}},
			// 2.2. Pre-consensus
			{NonConsensus: &mockNonConsensusMessage{Slot: 65, Type: ssv.SelectionProofPartialSig}},
			// 2.3. Consensus
			{Consensus: &mockConsensusMessage{Height: 101}},
			// 2.4. Post-consensus
			{NonConsensus: &mockNonConsensusMessage{Slot: 65, Type: ssv.PostConsensusPartialSig}},

			// 3. Lower height/slot:
			// 3.1 Decided
			{Consensus: &mockConsensusMessage{Height: 99, Decided: true}},
			// 3.2. Commit
			{Consensus: &mockConsensusMessage{Height: 99, Type: qbft.CommitMsgType}},
			// 3.3. Pre-consensus
			{NonConsensus: &mockNonConsensusMessage{Slot: 63, Type: ssv.SelectionProofPartialSig}},
		},
	}
}
