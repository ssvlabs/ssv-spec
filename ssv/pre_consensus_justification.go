package ssv

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

func (b *BaseRunner) shouldProcessingJustificationsForHeight(msg *qbft.SignedMessage) bool {
	rightQBFTHeight := b.QBFTController.CanStartInstance() == nil && b.QBFTController.Height+1 == msg.Message.Height
	hasData := len(msg.FullData) > 0
	return rightQBFTHeight && hasData
}

func (b *BaseRunner) validatePreConsensusJustifications(data *types.ConsensusData) error {
	if err := data.Validate(); err != nil {
		return err
	}

	// validate justification quorum
	if !b.Share.HasQuorum(len(data.PreConsensusJustifications)) {
		return errors.New("no quorum")
	}

	signers := make(map[types.OperatorID]bool)
	roots := make(map[[32]byte]bool)
	for i, msg := range data.PreConsensusJustifications {
		if err := msg.Validate(); err != nil {
			return err
		}

		// check unique signers
		if !signers[msg.Signer] {
			signers[msg.Signer] = true
		} else {
			return errors.New("duplicate signer")
		}

		// validate roots
		for _, msgRoot := range msg.Message.Messages {
			// validate roots
			if i == 0 {
				// record roots
				if !roots[msgRoot.SigningRoot] {
					roots[msgRoot.SigningRoot] = true
				}
			} else {
				// compare roots
				if !roots[msgRoot.SigningRoot] {
					return errors.New("invalid roots")
				}
			}
		}

		// verify sigs and duty.slot == msg.slot
		if err := b.validatePartialSigMsgForSlot(msg, data.Duty.Slot); err != nil {
			return err
		}
	}
	return nil
}

// processPreConsensusJustification processes pre-consensus justification
// highestDecidedDutySlot is the highest decided duty slot known
// is the qbft message carrying  the pre-consensus justification
func (b *BaseRunner) processPreConsensusJustification(runner Runner, highestDecidedDutySlot phase0.Slot, msg *qbft.SignedMessage) error {
	// TODO should validate qbft message?
	/**
	0) needs to process justifications
	1) validate message
		1.1) validate consensus data
		1.2) validate each signed msg
		1.3) validate quorum for justifications
		1.4) validate unique signers
		1.5) validate duty.slot == message slot
		1.6) validate message roots equal
		1.7) validate sigs
	2) if cd.Duty.Slot > highestDecidedDutySlot return nil
	3) if no running instance, run instance with consensus data duty
	4) add pre-consensus sigs to container
	5) decided on duty
	*/

	if !b.shouldProcessingJustificationsForHeight(msg) {
		return nil
	}

	cd := &types.ConsensusData{}
	if err := cd.Decode(msg.FullData); err != nil {
		return err
	}

	if highestDecidedDutySlot >= cd.Duty.Slot {
		return errors.New("duty.slot < highest decided slot")
	}

	if err := b.validatePreConsensusJustifications(cd); err != nil {
		return err
	}

	if !b.hasRunningDuty() {
		b.setupForNewDuty(&cd.Duty)
	}

	return nil
}
