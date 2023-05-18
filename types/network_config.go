package types

import (
	"math/big"
	"time"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// BeaconNetwork describes a network.
type BeaconNetwork struct {
	Name              string
	DefaultSyncOffset *big.Int // prod contract genesis block
	ForkVersion       [4]byte
	MinGenesisTime    uint64
	SlotDuration      time.Duration
	SlotsPerEpoch     uint64
	CapellaForkEpoch  uint64
}

var BeaconTestNetwork = BeaconNetwork{
	Name:              "now_test_network",
	DefaultSyncOffset: new(big.Int).SetInt64(8661727),
	ForkVersion:       [4]byte{0x99, 0x99, 0x99, 0x99},
	MinGenesisTime:    1616508000,
	SlotDuration:      12 * time.Second,
	SlotsPerEpoch:     32,
	CapellaForkEpoch:  162304,
}

func GetBeaconTestNetwork() BeaconNetwork {
	return BeaconTestNetwork
}

// EstimatedCurrentSlot returns the estimation of the current slot
func (n BeaconNetwork) EstimatedCurrentSlot() spec.Slot {
	return n.EstimatedSlotAtTime(time.Now().Unix())
}

// EstimatedSlotAtTime estimates slot at the given time
func (n BeaconNetwork) EstimatedSlotAtTime(time int64) spec.Slot {
	genesis := int64(n.MinGenesisTime)
	if time < genesis {
		return 0
	}
	return spec.Slot(uint64(time-genesis) / uint64(n.SlotDuration.Seconds()))
}

func (n BeaconNetwork) EstimatedTimeAtSlot(slot spec.Slot) int64 {
	d := int64(slot) * int64(n.SlotDuration.Seconds())
	return int64(n.MinGenesisTime) + d
}

// EstimatedCurrentEpoch estimates the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/beacon-chain.md#compute_start_slot_at_epoch
func (n BeaconNetwork) EstimatedCurrentEpoch() spec.Epoch {
	return n.EstimatedEpochAtSlot(n.EstimatedCurrentSlot())
}

// EstimatedEpochAtSlot estimates epoch at the given slot
func (n BeaconNetwork) EstimatedEpochAtSlot(slot spec.Slot) spec.Epoch {
	return spec.Epoch(slot / spec.Slot(n.SlotsPerEpoch))
}

func (n BeaconNetwork) FirstSlotAtEpoch(epoch spec.Epoch) spec.Slot {
	return spec.Slot(uint64(epoch) * n.SlotsPerEpoch)
}

func (n BeaconNetwork) EpochStartTime(epoch spec.Epoch) time.Time {
	firstSlot := n.FirstSlotAtEpoch(epoch)
	t := n.EstimatedTimeAtSlot(firstSlot)
	return time.Unix(t, 0)
}
