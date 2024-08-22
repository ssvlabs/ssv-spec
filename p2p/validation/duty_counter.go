package validation

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

// Holds the amount of duties for each identifier (validator/committee, duty role) for each epoch
type DutyCounter struct {
	duties map[types.MessageID]map[phase0.Epoch]map[phase0.Slot]struct{}
}

func NewDutyCounter() *DutyCounter {
	return &DutyCounter{
		duties: make(map[types.MessageID]map[phase0.Epoch]map[phase0.Slot]struct{}),
	}
}

// Records a duty
func (dc *DutyCounter) RecordDuty(msgID types.MessageID, epoch phase0.Epoch, slot phase0.Slot) {
	if _, exists := dc.duties[msgID]; !exists {
		dc.duties[msgID] = make(map[phase0.Epoch]map[phase0.Slot]struct{})
	}
	if _, exists := dc.duties[msgID][epoch]; !exists {
		dc.duties[msgID][epoch] = make(map[phase0.Slot]struct{})
	}
	dc.duties[msgID][epoch][slot] = struct{}{}
}

// Counts the number of duties for a given identifier and epoch
func (dc *DutyCounter) CountDutiesForEpoch(msgID types.MessageID, epoch phase0.Epoch) int {
	if _, exists := dc.duties[msgID]; !exists {
		return 0
	}
	if _, exists := dc.duties[msgID][epoch]; !exists {
		return 0
	}
	return len(dc.duties[msgID][epoch])
}

// Checks if has a duty
func (dc *DutyCounter) HasDuty(msgID types.MessageID, epoch phase0.Epoch, slot phase0.Slot) bool {
	if _, exists := dc.duties[msgID]; !exists {
		return false
	}
	if _, exists := dc.duties[msgID][epoch]; !exists {
		return false
	}
	_, exists := dc.duties[msgID][epoch][slot]
	return exists
}
