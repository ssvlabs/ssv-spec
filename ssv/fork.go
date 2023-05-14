package ssv

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// forkBasedOnLatestDecided will return domain type based on b.Height instance that was previously decided
func (b *BaseRunner) forkBasedOnLatestDecided() (types.DomainType, error) {
	inst := b.QBFTController.InstanceForHeight(b.QBFTController.Height)
	if inst == nil {
		return b.Share.DomainType, nil
	}

	_, decidedValue := inst.IsDecided()
	cd := &types.ConsensusData{}
	if err := cd.Decode(decidedValue); err != nil {
		return b.Share.DomainType, errors.Wrap(err, "could not decoded consensus data")
	}

	currentForkDigest := b.Share.DomainType
	for _, forkData := range currentForkDigest.GetForksData() {
		if b.BeaconNetwork.EstimatedEpochAtSlot(cd.Duty.Slot) >= forkData.Epoch {
			currentForkDigest = forkData.Domain
		}
	}
	return currentForkDigest, nil
}
