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
	// ForkEpochGloas is a test-only value placed after Fulu (arbitrary, like the other testnet
	// fork epochs here); adjust to match the node's config if cross-repo slot parity is wanted.
	ForkEpochGloas = 250000

	ForkEpochFulu = 242000

	ForkEpochPraterElectra = 232000

	ForkEpochPraterDeneb = 231680

	// ForkEpochPraterCapella Goerli taken from https://github.com/ethereum/execution-specs/blob/37a8f892341eb000e56e962a051a87e05a2e4443/network-upgrades/mainnet-upgrades/shanghai.md?plain=1#L18
	ForkEpochPraterCapella = 162304

	ForkEpochBellatrix = 144896

	ForkEpochAltair = 74240

	ForkEpochPhase0 = 0

	TestingDutyEpochPhase0         = TestingDutyEpoch
	TestingDutySlotPhase0          = TestingDutySlot
	TestingDutySlotPhase0NextEpoch = TestingDutySlot2
	TestingDutySlotPhase0Invalid   = TestingDutySlotPhase0 + 50

	TestingDutyEpochAltair         = ForkEpochAltair
	TestingDutySlotAltair          = ForkEpochAltair * 32
	TestingDutySlotAltairNextEpoch = TestingDutySlotAltair + 32
	TestingDutySlotAltairInvalid   = TestingDutySlotAltair + 50

	TestingDutyEpochBellatrix         = ForkEpochBellatrix
	TestingDutySlotBellatrix          = ForkEpochBellatrix * 32
	TestingDutySlotBellatrixNextEpoch = TestingDutySlotBellatrix + 32
	TestingDutySlotBellatrixInvalid   = TestingDutySlotBellatrix + 50

	TestingDutyEpochCapella         = ForkEpochPraterCapella
	TestingDutySlotCapella          = ForkEpochPraterCapella * 32
	TestingDutySlotCapellaNextEpoch = TestingDutySlotCapella + 32
	TestingDutySlotCapellaInvalid   = TestingDutySlotCapella + 50

	TestingDutyEpochDeneb         = ForkEpochPraterDeneb
	TestingDutySlotDeneb          = ForkEpochPraterDeneb * 32
	TestingDutySlotDenebNextEpoch = TestingDutySlotDeneb + 32
	TestingDutySlotDenebInvalid   = TestingDutySlotDeneb + 50

	TestingDutyEpochElectra         = ForkEpochPraterElectra
	TestingDutySlotElectra          = ForkEpochPraterElectra*32 + 12
	TestingDutySlotElectraNextEpoch = TestingDutySlotElectra + 32
	TestingDutySlotElectraInvalid   = TestingDutySlotElectra + 50

	TestingDutyEpochFulu         = ForkEpochFulu
	TestingDutySlotFulu          = ForkEpochFulu*32 + 12
	TestingDutySlotFuluNextEpoch = TestingDutySlotFulu + 32
	TestingDutySlotFuluInvalid   = TestingDutySlotFulu + 50

	TestingDutyEpochGloas         = ForkEpochGloas
	TestingDutySlotGloas          = ForkEpochGloas*32 + 12
	TestingDutySlotGloasNextEpoch = TestingDutySlotGloas + 32
	TestingDutySlotGloasInvalid   = TestingDutySlotGloas + 50
)

var TestingDutyEpochV = func(version spec.DataVersion) phase0.Epoch {
	switch version {
	case spec.DataVersionPhase0:
		return TestingDutyEpochPhase0
	case spec.DataVersionAltair:
		return TestingDutyEpochAltair
	case spec.DataVersionBellatrix:
		return TestingDutyEpochBellatrix
	case spec.DataVersionCapella:
		return TestingDutyEpochCapella
	case spec.DataVersionDeneb:
		return TestingDutyEpochDeneb
	case spec.DataVersionElectra:
		return TestingDutyEpochElectra
	case spec.DataVersionFulu:
		return TestingDutyEpochFulu
	case gloas.DataVersionGloas:
		return TestingDutyEpochGloas

	default:
		panic("unsupported version")
	}
}

var TestingDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionPhase0:
		return TestingDutySlotPhase0
	case spec.DataVersionAltair:
		return TestingDutySlotAltair
	case spec.DataVersionBellatrix:
		return TestingDutySlotBellatrix
	case spec.DataVersionCapella:
		return TestingDutySlotCapella
	case spec.DataVersionDeneb:
		return TestingDutySlotDeneb
	case spec.DataVersionElectra:
		return TestingDutySlotElectra
	case spec.DataVersionFulu:
		return TestingDutySlotFulu
	case gloas.DataVersionGloas:
		return TestingDutySlotGloas

	default:
		panic("unsupported version")
	}
}

var TestingDutySlotNextEpochV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionPhase0:
		return TestingDutySlotPhase0NextEpoch
	case spec.DataVersionAltair:
		return TestingDutySlotAltairNextEpoch
	case spec.DataVersionBellatrix:
		return TestingDutySlotBellatrixNextEpoch
	case spec.DataVersionCapella:
		return TestingDutySlotCapellaNextEpoch
	case spec.DataVersionDeneb:
		return TestingDutySlotDenebNextEpoch
	case spec.DataVersionElectra:
		return TestingDutySlotElectraNextEpoch
	case spec.DataVersionFulu:
		return TestingDutySlotFuluNextEpoch
	case gloas.DataVersionGloas:
		return TestingDutySlotGloasNextEpoch

	default:
		panic("unsupported version")
	}
}

var TestingInvalidDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionPhase0:
		return TestingDutySlotPhase0Invalid
	case spec.DataVersionAltair:
		return TestingDutySlotAltairInvalid
	case spec.DataVersionBellatrix:
		return TestingDutySlotBellatrixInvalid
	case spec.DataVersionCapella:
		return TestingDutySlotCapellaInvalid
	case spec.DataVersionDeneb:
		return TestingDutySlotDenebInvalid
	case spec.DataVersionElectra:
		return TestingDutySlotElectraInvalid
	case spec.DataVersionFulu:
		return TestingDutySlotFuluInvalid
	case gloas.DataVersionGloas:
		return TestingDutySlotGloasInvalid

	default:
		panic("unsupported version")
	}
}

// VersionBySlot resolves a slot's fork. slot < ForkEpoch*32 is equivalent to slot/32 < ForkEpoch, so
// it delegates rather than repeating every threshold: one table to keep in step with the fork schedule.
var VersionBySlot = func(slot phase0.Slot) spec.DataVersion {
	return VersionByEpoch(phase0.Epoch(slot / 32))
}

var VersionByEpoch = func(epoch phase0.Epoch) spec.DataVersion {
	if epoch < ForkEpochAltair {
		return spec.DataVersionPhase0
	} else if epoch < ForkEpochBellatrix {
		return spec.DataVersionAltair
	} else if epoch < ForkEpochPraterCapella {
		return spec.DataVersionBellatrix
	} else if epoch < ForkEpochPraterDeneb {
		return spec.DataVersionCapella
	} else if epoch < ForkEpochPraterElectra {
		return spec.DataVersionDeneb
	} else if epoch < ForkEpochFulu {
		return spec.DataVersionElectra
	} else if epoch < ForkEpochGloas {
		return spec.DataVersionFulu
	}
	return gloas.DataVersionGloas
}
