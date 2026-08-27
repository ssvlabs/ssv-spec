package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ProposedBlockRoots records each slot's §4-decided Gloas block root — the §4→§6 linkage (SIP #94 §6):
// the proposer runner records the root of the block its cluster decided, and the envelope runner and
// its value check read it, tying the §6 envelope to exactly that block. Production shares one instance
// between the two runners; pruning is node-side.
type ProposedBlockRoots map[phase0.Slot]phase0.Root

// Record stores the slot's decided block root
func (roots ProposedBlockRoots) Record(slot phase0.Slot, root phase0.Root) {
	roots[slot] = root
}

// Get returns the slot's decided block root, if recorded
func (roots ProposedBlockRoots) Get(slot phase0.Slot) (phase0.Root, bool) {
	root, ok := roots[slot]
	return root, ok
}
