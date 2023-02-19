package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
)

// processPreConsensusJustification processes pre-consensus justification
// highestDecidedDutySlot is the highest decided duty slot known
// is the qbft message carrying  the pre-consensus justification
func (b *BaseRunner) processPreConsensusJustification(runner Runner, highestDecidedDutySlot phase0.Slot, msg *qbft.SignedMessage) error {
	/**
	1) validate message
		1.1) validate consensus data
		1.2) validate duty.slot == message slot
		1.3) validate unique signers
		1.4) validate quorum for justifications
		1.5) validate message roots equal
		1.6) validate sigs
		1.7)
	2) if cd.Duty.Slot > highestDecidedDutySlot return nil
	3) if no running instance, run instance with consensus data duty
	4) add pre-consensus sigs to container
	5) decided on duty
	*/
	return nil
}
