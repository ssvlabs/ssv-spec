package queue

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/queue"
)

func AllInOne() *MessagePriorityTest {
	return &MessagePriorityTest{
		Name: "All-in-one",
		State: &queue.State{
			HasRunningInstance: true,
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
			// 1.1.1. Consensus/Prepare
			mockConsensusMessage{Height: 100, Type: qbft.PrepareMsgType},
			// 1.1.2. Consensus/Proposal
			mockConsensusMessage{Height: 100, Type: qbft.ProposalMsgType},
			// 1.1.3. Consensus/Commit
			mockConsensusMessage{Height: 100, Type: qbft.CommitMsgType},
			// 1.1.4. Consensus/<Other>
			mockConsensusMessage{Height: 100, Type: qbft.RoundChangeMsgType},
			// 1.2. Pre-consensus
			mockNonConsensusMessage{Slot: 64, Type: ssv.SelectionProofPartialSig},
			// 1.3. Post-consensus
			mockNonConsensusMessage{Slot: 64, Type: ssv.PostConsensusPartialSig},

			// 2. Higher height/slot:
			// 2.1 Decided
			mockConsensusMessage{Height: 101, Decided: true},
			// 2.2. Pre-consensus
			mockNonConsensusMessage{Slot: 65, Type: ssv.SelectionProofPartialSig},
			// 2.3. Consensus
			mockConsensusMessage{Height: 101},
			// 2.4. Post-consensus
			mockNonConsensusMessage{Slot: 65, Type: ssv.PostConsensusPartialSig},

			// 3. Lower height/slot:
			// 3.1 Decided
			mockConsensusMessage{Height: 99, Decided: true},
			// 3.2. Commit
			mockConsensusMessage{Height: 99, Type: qbft.CommitMsgType},
			// 3.3. Pre-consensus
			mockNonConsensusMessage{Slot: 63, Type: ssv.SelectionProofPartialSig},
		},
	}
}
