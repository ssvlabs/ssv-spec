package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ==================================================
// Beacon Fork Epochs and Slots (Main, Next, Invalid)
// ==================================================

const (

	// Electra Fork Epoch: TODO - update to the correct value
	ForkEpochPraterElectra = 500000

	//Deneb Fork Epoch
	ForkEpochPraterDeneb = 231680

	// ForkEpochPraterCapella Goerli taken from https://github.com/ethereum/execution-specs/blob/37a8f892341eb000e56e962a051a87e05a2e4443/network-upgrades/mainnet-upgrades/shanghai.md?plain=1#L18
	ForkEpochPraterCapella = 162304

	TestingDutyEpochCapella         = ForkEpochPraterCapella
	TestingDutySlotCapella          = ForkEpochPraterCapella * 32
	TestingDutySlotCapellaNextEpoch = TestingDutySlotCapella + 32
	TestingDutySlotCapellaInvalid   = TestingDutySlotCapella + 50

	TestingDutyEpochDeneb         = ForkEpochPraterDeneb
	TestingDutySlotDeneb          = ForkEpochPraterDeneb * 32
	TestingDutySlotDenebNextEpoch = TestingDutySlotDeneb + 32
	TestingDutySlotDenebInvalid   = TestingDutySlotDeneb + 50

	TestingDutyEpochElectra         = ForkEpochPraterElectra
	TestingDutySlotElectra          = ForkEpochPraterElectra * 32
	TestingDutySlotElectraNextEpoch = TestingDutySlotElectra + 32
	TestingDutySlotElectraInvalid   = TestingDutySlotElectra + 50
)

var TestingDutyEpochV = func(version spec.DataVersion) phase0.Epoch {
	switch version {
	case spec.DataVersionCapella:
		return TestingDutyEpochCapella
	case spec.DataVersionDeneb:
		return TestingDutyEpochDeneb
	case spec.DataVersionElectra:
		return TestingDutyEpochElectra

	default:
		panic("unsupported version")
	}
}

var TestingDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionCapella:
		return TestingDutySlotCapella
	case spec.DataVersionDeneb:
		return TestingDutySlotDeneb
	case spec.DataVersionElectra:
		return TestingDutySlotElectra

	default:
		panic("unsupported version")
	}
}

var TestingDutySlotNextEpochV = func(version spec.DataVersion) phase0.Slot {

	switch version {
	case spec.DataVersionCapella:
		return TestingDutySlotCapellaNextEpoch
	case spec.DataVersionDeneb:
		return TestingDutySlotDenebNextEpoch
	case spec.DataVersionElectra:
		return TestingDutySlotElectraNextEpoch

	default:
		panic("unsupported version")
	}
}

var TestingInvalidDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionCapella:
		return TestingDutySlotCapellaInvalid
	case spec.DataVersionDeneb:
		return TestingDutySlotDenebInvalid
	case spec.DataVersionElectra:
		return TestingDutySlotElectraInvalid

	default:
		panic("unsupported version")
	}
}

var VersionBySlot = func(slot phase0.Slot) spec.DataVersion {
	if slot < ForkEpochPraterDeneb*32 {
		return spec.DataVersionCapella
	} else if slot < ForkEpochPraterElectra*32 {
		return spec.DataVersionDeneb
	}
	return spec.DataVersionElectra
}

var VersionByEpoch = func(epoch phase0.Epoch) spec.DataVersion {
	if epoch < ForkEpochPraterDeneb {
		return spec.DataVersionCapella
	} else if epoch < ForkEpochPraterElectra {
		return spec.DataVersionDeneb
	}
	return spec.DataVersionElectra
}
