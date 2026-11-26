package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// Beacon Fork Epochs and Slots (Main, Next, Invalid)
// ==================================================

const (
	// ForkEpochGloas is a test-only value (arbitrary, like a testnet fork epoch); adjust to match the node's
	// config if cross-repo slot parity is wanted.
	ForkEpochGloas = 250000

	TestingDutyEpochGloas         = ForkEpochGloas
	TestingDutySlotGloas          = ForkEpochGloas*32 + 12
	TestingDutySlotGloasNextEpoch = TestingDutySlotGloas + 32
	TestingDutySlotGloasInvalid   = TestingDutySlotGloas + 50

	// TestingEpochDependentFrom is the first epoch for which the testing beacon node has a fork version
	// (TestingForkVersion) and a proposer-duties dependent root (TestingProposerDutiesDependentRootAt) of its
	// own. The epochs before it, where the other fixtures' slots lie, share one of each.
	TestingEpochDependentFrom = ForkEpochGloas + 100
	// TestingDutySlotEpochDependent is a duty slot past TestingEpochDependentFrom. A duty there pins the epoch
	// its runner asks the beacon node for: the answer for any other epoch gives another signing root.
	TestingDutySlotEpochDependent = (TestingEpochDependentFrom+1)*32 + 12
)

// TestingForkVersion is the fork version the testing beacon node computes an epoch's signing domains with.
func TestingForkVersion(epoch phase0.Epoch) phase0.Version {
	if epoch < TestingEpochDependentFrom {
		return types.GenesisForkVersion
	}
	n := uint64(epoch-TestingEpochDependentFrom) + 1
	return phase0.Version{0xfe, byte(n), byte(n >> 8), byte(n >> 16)}
}

var TestingDutyEpochV = func(version spec.DataVersion) phase0.Epoch {
	if version != gloas.DataVersionGloas {
		panic("unsupported version")
	}
	return TestingDutyEpochGloas
}

var TestingDutySlotV = func(version spec.DataVersion) phase0.Slot {
	if version != gloas.DataVersionGloas {
		panic("unsupported version")
	}
	return TestingDutySlotGloas
}

var TestingDutySlotNextEpochV = func(version spec.DataVersion) phase0.Slot {
	if version != gloas.DataVersionGloas {
		panic("unsupported version")
	}
	return TestingDutySlotGloasNextEpoch
}

var TestingInvalidDutySlotV = func(version spec.DataVersion) phase0.Slot {
	if version != gloas.DataVersionGloas {
		panic("unsupported version")
	}
	return TestingDutySlotGloasInvalid
}
