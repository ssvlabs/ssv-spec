package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ProposedBlock records the facts of a slot's §4-decided Gloas block that the §6 envelope duty binds a
// disseminated envelope against (SIP #94 §6): the block root, the block's parent root, and the bid's
// execution-requests root. The proposer runner extracts all three from the decided block; the envelope
// runner reads them to content-select the envelope that matches the cluster's own §4 decision.
type ProposedBlock struct {
	BlockRoot             phase0.Root
	ParentRoot            phase0.Root
	ExecutionRequestsRoot phase0.Root
}

// ProposedBlocks is the §4→§6 linkage store: slot → decided-block facts. Production shares one instance
// between the proposer runner (writer) and the envelope runner (reader); pruning is node-side.
type ProposedBlocks map[phase0.Slot]ProposedBlock

// Record stores the slot's decided-block facts.
func (b ProposedBlocks) Record(slot phase0.Slot, block ProposedBlock) {
	b[slot] = block
}

// Get returns the slot's decided-block facts, if recorded.
func (b ProposedBlocks) Get(slot phase0.Slot) (ProposedBlock, bool) {
	block, ok := b[slot]
	return block, ok
}
