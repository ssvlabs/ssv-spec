package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

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
)

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
