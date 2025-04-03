package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/types"
)

var SupportedAggregatorVersions = []spec.DataVersion{spec.DataVersionPhase0, spec.DataVersionElectra}

// ==================================================
// Versioned Aggregator Duty
// ==================================================

var TestingAggregatorDuty = func(version spec.DataVersion) *types.ValidatorDuty {
	return &types.ValidatorDuty{
		Type:                    types.BNRoleAggregator,
		PubKey:                  TestingValidatorPubKey,
		Slot:                    TestingDutySlotV(version),
		ValidatorIndex:          TestingValidatorIndex,
		CommitteeIndex:          22,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}
}

var TestingAggregatorDutyNextEpoch = func(version spec.DataVersion) *types.ValidatorDuty {
	return &types.ValidatorDuty{
		Type:                    types.BNRoleAggregator,
		PubKey:                  TestingValidatorPubKey,
		Slot:                    TestingDutySlotNextEpochV(version),
		ValidatorIndex:          TestingValidatorIndex,
		CommitteeIndex:          22,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}
}

var TestingAggregatorDutyFirstSlot = func() *types.ValidatorDuty {
	return &types.ValidatorDuty{
		Type:                    types.BNRoleAggregator,
		PubKey:                  TestingValidatorPubKey,
		Slot:                    0,
		ValidatorIndex:          TestingValidatorIndex,
		CommitteeIndex:          22,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}
}

// ==================================================
// Versioned AggregateAndProof
// ==================================================

var TestingAggregateAndProofV = func(version spec.DataVersion) ssz.Marshaler {
	if version == spec.DataVersionElectra {
		return TestingElectraAggregateAndProof
	} else {
		return TestingPhase0AggregateAndProof
	}
}

var TestingVersionedSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) *spec.VersionedSignedAggregateAndProof {

	switch version {
	case spec.DataVersionPhase0:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Phase0:  TestingPhase0SignedAggregateAndProof(ks),
		}

	case spec.DataVersionAltair:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Altair:  TestingPhase0SignedAggregateAndProof(ks),
		}

	case spec.DataVersionBellatrix:
		return &spec.VersionedSignedAggregateAndProof{
			Version:   version,
			Bellatrix: TestingPhase0SignedAggregateAndProof(ks),
		}

	case spec.DataVersionCapella:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Capella: TestingPhase0SignedAggregateAndProof(ks),
		}
	case spec.DataVersionDeneb:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Deneb:   TestingPhase0SignedAggregateAndProof(ks),
		}
	case spec.DataVersionElectra:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Electra: TestingElectraSignedAggregateAndProof(ks),
		}
	default:
		panic("unknown data version")
	}
}

var TestingSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	switch version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix, spec.DataVersionCapella, spec.DataVersionDeneb:
		return TestingPhase0SignedAggregateAndProof(ks)
	case spec.DataVersionElectra:
		return TestingElectraSignedAggregateAndProof(ks)
	default:
		panic("unknown data version")
	}
}

var TestingAggregateAndProofBytesV = func(version spec.DataVersion) []byte {
	if version == spec.DataVersionElectra {
		return TestingElectraAggregateAndProofBytes
	} else {
		return TestingPhase0AggregateAndProofBytes
	}
}

var TestingWrongAggregateAndProofV = func(version spec.DataVersion) ssz.Marshaler {
	if version == spec.DataVersionElectra {
		return TestingWrongElectraAggregateAndProof
	} else {
		return TestingWrongPhase0AggregateAndProof
	}
}

// phase0.AggregateAndProof

var TestingPhase0AggregateAndProof = &phase0.AggregateAndProof{
	AggregatorIndex: 1,
	SelectionProof:  phase0.BLSSignature{},
	Aggregate: &phase0.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Signature:       phase0.BLSSignature{},
		Data:            TestingAttestationData(spec.DataVersionPhase0),
	},
}
var TestingPhase0AggregateAndProofBytes = func() []byte {
	ret, _ := TestingPhase0AggregateAndProof.MarshalSSZ()
	return ret
}()

var TestingWrongPhase0AggregateAndProof = func() *phase0.AggregateAndProof {
	byts, err := TestingPhase0AggregateAndProof.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	ret := &phase0.AggregateAndProof{}
	if err := ret.UnmarshalSSZ(byts); err != nil {
		panic(err.Error())
	}
	ret.AggregatorIndex = 100
	return ret
}()

// electra.AggregateAndProof

var TestingElectraAggregateAndProof = &electra.AggregateAndProof{
	AggregatorIndex: 1,
	SelectionProof:  phase0.BLSSignature{},
	Aggregate: &electra.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Signature:       phase0.BLSSignature{},
		Data:            TestingAttestationData(spec.DataVersionElectra),
		CommitteeBits:   bitfield.NewBitvector64(),
	},
}
var TestingElectraAggregateAndProofBytes = func() []byte {
	ret, _ := TestingElectraAggregateAndProof.MarshalSSZ()
	return ret
}()
var TestingWrongElectraAggregateAndProof = func() *electra.AggregateAndProof {
	byts, err := TestingElectraAggregateAndProof.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	ret := &electra.AggregateAndProof{}
	if err := ret.UnmarshalSSZ(byts); err != nil {
		panic(err.Error())
	}
	ret.AggregatorIndex = 100
	return ret
}()

var TestingPhase0SignedAggregateAndProof = func(ks *TestKeySet) *phase0.SignedAggregateAndProof {
	return &phase0.SignedAggregateAndProof{
		Message:   TestingPhase0AggregateAndProof,
		Signature: signBeaconObject(TestingPhase0AggregateAndProof, types.DomainAggregateAndProof, ks),
	}
}

var TestingElectraSignedAggregateAndProof = func(ks *TestKeySet) *electra.SignedAggregateAndProof {
	return &electra.SignedAggregateAndProof{
		Message:   TestingElectraAggregateAndProof,
		Signature: signBeaconObject(TestingElectraAggregateAndProof, types.DomainAggregateAndProof, ks),
	}
}
