package testingutilscomparable

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/ssv-spec/types"
)

func compareSyncCommitteeIndices(source, target []uint64) []Difference {
	ret := make([]Difference, 0)

	if len(source) != len(target) {
		ret = append(ret, Sprintf("Length", "source %d != target %d", len(source), len(target)))
		return ret
	}

	for i := range source {
		if source[i] != target[i] {
			ret = append(ret, Sprintf(fmt.Sprintf("Item %d", i), "source %s != target %s", source[i], target[i]))
		}
	}

	return ret
}

func CompareDuty(source, target *types.Duty) []Difference {
	ret := make([]Difference, 0)

	if (source == nil && target != nil) || (source != nil && target == nil) {
		ret = append(ret, Sprintf("Duty nil?", "source %t != target %t", source == nil, target == nil))
		return ret
	}

	if source.Type != target.Type {
		ret = append(ret, Sprintf("Type", "source %s != target %s", source.Type.String(), target.Type.String()))
	}

	if !bytes.Equal(source.PubKey[:], target.PubKey[:]) {
		ret = append(ret, Sprintf("PubKey", "source %x != target %x", source.PubKey, target.PubKey))
	}

	if source.Slot != target.Slot {
		ret = append(ret, Sprintf("Slot", "source %d != target %d", source.Slot, target.Slot))
	}

	if source.ValidatorIndex != target.ValidatorIndex {
		ret = append(ret, Sprintf("ValidatorIndex", "source %d != target %d", source.ValidatorIndex, target.ValidatorIndex))
	}

	if source.CommitteeIndex != target.CommitteeIndex {
		ret = append(ret, Sprintf("CommitteeIndex", "source %d != target %d", source.CommitteeIndex, target.CommitteeIndex))
	}

	if source.CommitteeLength != target.CommitteeLength {
		ret = append(ret, Sprintf("CommitteeLength", "source %d != target %d", source.CommitteeLength, target.CommitteeLength))
	}

	if source.CommitteesAtSlot != target.CommitteesAtSlot {
		ret = append(ret, Sprintf("CommitteesAtSlot", "source %d != target %d", source.CommitteesAtSlot, target.CommitteesAtSlot))
	}

	if source.ValidatorCommitteeIndex != target.ValidatorCommitteeIndex {
		ret = append(ret, Sprintf("ValidatorCommitteeIndex", "source %d != target %d", source.ValidatorCommitteeIndex, target.ValidatorCommitteeIndex))
	}

	if diff := NestedCompare("ValidatorSyncCommitteeIndices", compareSyncCommitteeIndices(source.ValidatorSyncCommitteeIndices, target.ValidatorSyncCommitteeIndices)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	return ret
}
