package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	ssz "github.com/ferranbt/fastssz"

	"github.com/bloxapp/ssv-spec/types"
)

var TestingBeaconBlockV = func(version spec.DataVersion) ssz.HashRoot {
	switch version {
	case spec.DataVersionBellatrix:
		return TestingBeaconBlock
	default:
		panic("unsupported version")
	}
}

var TestingBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	tBlk := TestingBeaconBlockV(version)
	var ret []byte

	switch version {
	case spec.DataVersionBellatrix:
		blk, ok := tBlk.(*bellatrix.BeaconBlock)
		if !ok {
			panic("failed to cast")
		}
		ret, _ = blk.MarshalSSZ()

	default:
		panic("unsupported version")
	}

	return ret
}

var TestingBlindedBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	switch version {
	case spec.DataVersionBellatrix:
		ret, _ := TestingBlindedBeaconBlock.MarshalSSZ()
		return ret

	default:
		panic("unsupported version")
	}
}

var TestingWrongBeaconBlockV = func(version spec.DataVersion) ssz.HashRoot {
	blkByts := TestingBeaconBlockBytesV(version)

	switch version {
	case spec.DataVersionBellatrix:
		ret := &bellatrix.BeaconBlock{}
		if err := ret.UnmarshalSSZ(blkByts); err != nil {
			panic(err.Error())
		}
		ret.Slot = 100
		return ret

	default:
		panic("unsupported version")
	}
}

var TestingSignedBeaconBlockV = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	tBlk := TestingBeaconBlockV(version)

	switch version {
	case spec.DataVersionBellatrix:
		blk, ok := tBlk.(*bellatrix.BeaconBlock)
		if !ok {
			panic("failed to cast")
		}
		return &bellatrix.SignedBeaconBlock{
			Message:   blk,
			Signature: signBeaconObject(blk, types.DomainProposer, ks),
		}

	default:
		panic("unsupported version")
	}
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
		duty.Slot = TestingDutySlot

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
		duty.Slot = TestingDutySlot2

	default:
		panic("unsupported version")
	}

	return duty
}
