package testingutils

import (
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/bloxapp/ssv-spec/types"
)

const (
	// ForkSlotCapella taken from https://github.com/ethereum/consensus-specs/blob/1c424d76eddbacae3cbffed8276264b46951456b/specs/capella/fork.md?plain=1#L30
	ForkSlotCapella = 6209536 // Epoch(194048)
	// TestingDutySlotBellatrix keeping this value to not break the test roots
	TestingDutySlotBellatrix          = 12
	TestingDutySlotBellatrixNextEpoch = 50
	TestingDutySlotBellatrixInvalid   = 50
)

var TestingBeaconBlockV = func(version spec.DataVersion) *spec.VersionedBeaconBlock {
	switch version {
	case spec.DataVersionBellatrix:
		return &spec.VersionedBeaconBlock{
			Version:   version,
			Bellatrix: TestingBeaconBlock,
		}
	default:
		panic("unsupported version")
	}
}

var TestingBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	var ret []byte
	vBlk := TestingBeaconBlockV(version)
	if vBlk.IsEmpty() {
		panic("empty block")
	}

	switch version {
	case spec.DataVersionBellatrix:
		ret, _ = vBlk.Bellatrix.MarshalSSZ()

	default:
		panic("unsupported version")
	}

	return ret
}

var TestingBlindedBeaconBlockV = func(version spec.DataVersion) *api.VersionedBlindedBeaconBlock {
	switch version {
	case spec.DataVersionBellatrix:
		return &api.VersionedBlindedBeaconBlock{
			Version:   version,
			Bellatrix: TestingBlindedBeaconBlock,
		}
	default:
		panic("unsupported version")
	}
}

var TestingBlindedBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	var ret []byte
	vBlk := TestingBlindedBeaconBlockV(version)
	if vBlk.IsEmpty() {
		panic("empty block")
	}

	switch version {
	case spec.DataVersionBellatrix:
		ret, _ = vBlk.Bellatrix.MarshalSSZ()

	default:
		panic("unsupported version")
	}

	return ret
}

var TestingWrongBeaconBlockV = func(version spec.DataVersion) *spec.VersionedBeaconBlock {
	blkByts := TestingBeaconBlockBytesV(version)

	switch version {
	case spec.DataVersionBellatrix:
		ret := &bellatrix.BeaconBlock{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Slot = 100
		return &spec.VersionedBeaconBlock{
			Version:   version,
			Bellatrix: ret,
		}

	default:
		panic("unsupported version")
	}
}

var TestingSignedBeaconBlockV = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	vBlk := TestingBeaconBlockV(version)
	if vBlk.IsEmpty() {
		panic("empty block")
	}

	switch version {
	case spec.DataVersionBellatrix:
		return &bellatrix.SignedBeaconBlock{
			Message:   vBlk.Bellatrix,
			Signature: signBeaconObject(vBlk.Bellatrix, types.DomainProposer, ks),
		}

	default:
		panic("unsupported version")
	}
}

var VersionBySlot = func(slot phase0.Slot) spec.DataVersion {
	if slot < ForkSlotCapella {
		return spec.DataVersionBellatrix
	}
	return spec.DataVersionCapella
}

var TestingProposerDutyV = func(version spec.DataVersion) *types.Duty {
	duty := &types.Duty{
		Type:                    types.BNRoleProposer,
		PubKey:                  TestingValidatorPubKey,
		ValidatorIndex:          TestingValidatorIndex,
		CommitteeIndex:          3,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}

	switch version {
	case spec.DataVersionBellatrix:
		duty.Slot = TestingDutySlotBellatrix

	default:
		panic("unsupported version")
	}

	return duty
}

var TestingProposerDutyNextEpochV = func(version spec.DataVersion) *types.Duty {
	duty := &types.Duty{
		Type:                    types.BNRoleProposer,
		PubKey:                  TestingValidatorPubKey,
		ValidatorIndex:          TestingValidatorIndex,
		CommitteeIndex:          3,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}

	switch version {
	case spec.DataVersionBellatrix:
		duty.Slot = TestingDutySlotBellatrixNextEpoch

	default:
		panic("unsupported version")
	}

	return duty
}

var TestingInvalidDutySlotV = func(version spec.DataVersion) phase0.Slot {
	switch version {
	case spec.DataVersionBellatrix:
		return TestingDutySlotBellatrixInvalid

	default:
		panic("unsupported version")
	}
}
