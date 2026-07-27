package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

var SupportedAggregatorVersions = []spec.DataVersion{spec.DataVersionPhase0, spec.DataVersionElectra, gloas.DataVersionGloas}

// ==================================================
// Versioned Aggregator Duty
// ==================================================

var TestingAggregatorDuty = func(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDutyOnlyAggregator(version)
}

var TestingAggregatorDutyNextEpoch = func(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	d := TestingAggregatorCommitteeDutyOnlyAggregator(version)
	d.Slot = TestingDutySlotNextEpochV(version)
	for i := range d.ValidatorDuties {
		d.ValidatorDuties[i].Slot = TestingDutySlotNextEpochV(version)
	}
	return d
}

var TestingAggregatorDutyFirstSlot = func() *types.AggregatorCommitteeDuty {
	d := TestingAggregatorCommitteeDutyOnlyAggregator(spec.DataVersionPhase0)
	d.Slot = 0
	for i := range d.ValidatorDuties {
		d.ValidatorDuties[i].Slot = 0
	}
	return d
}

// ==================================================
// Versioned AggregateAndProof
// ==================================================

var TestingAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	if version >= spec.DataVersionElectra {
		return TestingElectraAggregateAndProofV(aggregatorIndex, version)
	}
	return TestingPhase0AggregateAndProof(aggregatorIndex)
}

var TestingVersionedSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) *spec.VersionedSignedAggregateAndProof {

	switch version {
	case spec.DataVersionPhase0:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Phase0:  TestingPhase0SignedAggregateAndProof(ks, TestingValidatorIndex),
		}

	case spec.DataVersionAltair:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Altair:  TestingPhase0SignedAggregateAndProof(ks, TestingValidatorIndex),
		}

	case spec.DataVersionBellatrix:
		return &spec.VersionedSignedAggregateAndProof{
			Version:   version,
			Bellatrix: TestingPhase0SignedAggregateAndProof(ks, TestingValidatorIndex),
		}

	case spec.DataVersionCapella:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Capella: TestingPhase0SignedAggregateAndProof(ks, TestingValidatorIndex),
		}
	case spec.DataVersionDeneb:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Deneb:   TestingPhase0SignedAggregateAndProof(ks, TestingValidatorIndex),
		}
	case spec.DataVersionElectra:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Electra: TestingElectraSignedAggregateAndProofV(ks, TestingValidatorIndex, version),
		}
	case spec.DataVersionFulu:
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Fulu:    TestingElectraSignedAggregateAndProofV(ks, TestingValidatorIndex, version),
		}
	case gloas.DataVersionGloas:
		// Gloas reuses the Electra aggregate shape (SIP #94 §2).
		return &spec.VersionedSignedAggregateAndProof{
			Version: version,
			Electra: TestingElectraSignedAggregateAndProofV(ks, TestingValidatorIndex, version),
		}
	default:
		panic("unknown data version")
	}
}

var TestingSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	switch version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix, spec.DataVersionCapella, spec.DataVersionDeneb:
		return TestingPhase0SignedAggregateAndProof(ks, TestingValidatorIndex)
	case spec.DataVersionElectra, spec.DataVersionFulu, gloas.DataVersionGloas:
		return TestingElectraSignedAggregateAndProofV(ks, TestingValidatorIndex, version)
	default:
		panic("unknown data version")
	}
}

var TestingAggregateAndProofBytesV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) []byte {
	if version >= spec.DataVersionElectra {
		return TestingElectraAggregateAndProofBytesV(aggregatorIndex, version)
	}
	return TestingPhase0AggregateAndProofBytes(aggregatorIndex)
}

var TestingWrongAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	if version >= spec.DataVersionElectra {
		return TestingWrongElectraAggregateAndProofV(aggregatorIndex, version)
	}
	return TestingWrongPhase0AggregateAndProof(aggregatorIndex)
}

// phase0.AggregateAndProof

var TestingPhase0AggregateAndProof = func(aggregatorIndex phase0.ValidatorIndex) *phase0.AggregateAndProof {
	return &phase0.AggregateAndProof{
		AggregatorIndex: aggregatorIndex,
		SelectionProof:  phase0.BLSSignature{},
		Aggregate: &phase0.Attestation{
			AggregationBits: bitfield.NewBitlist(128),
			Signature:       phase0.BLSSignature{},
			Data:            TestingAttestationData(spec.DataVersionPhase0),
		},
	}
}
var TestingPhase0AggregateAndProofBytes = func(aggregatorIndex phase0.ValidatorIndex) []byte {
	ret, _ := TestingPhase0AggregateAndProof(aggregatorIndex).MarshalSSZ()
	return ret
}

var TestingWrongPhase0AggregateAndProof = func(aggregatorIndex phase0.ValidatorIndex) *phase0.AggregateAndProof {
	byts, err := TestingPhase0AggregateAndProof(aggregatorIndex).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	ret := &phase0.AggregateAndProof{}
	if err := ret.UnmarshalSSZ(byts); err != nil {
		panic(err.Error())
	}
	ret.AggregatorIndex = 100
	return ret
}

// electra.AggregateAndProof

// TestingElectraAggregateAndProofV builds the Electra-shaped AggregateAndProof for the given fork.
// The version matters from Gloas on: the attestation data carries the payload-status index and the
// Gloas duty slot (SIP #94 §2), so fixtures must match the data the runner aggregates over.
var TestingElectraAggregateAndProofV = func(aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) *electra.AggregateAndProof {
	return &electra.AggregateAndProof{
		AggregatorIndex: aggregatorIndex,
		SelectionProof:  phase0.BLSSignature{},
		Aggregate: &electra.Attestation{
			AggregationBits: bitfield.NewBitlist(128),
			Signature:       phase0.BLSSignature{},
			Data:            TestingAttestationData(version),
			CommitteeBits:   bitfield.NewBitvector64(),
		},
	}
}

var TestingElectraAggregateAndProofBytesV = func(aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) []byte {
	ret, _ := TestingElectraAggregateAndProofV(aggregatorIndex, version).MarshalSSZ()
	return ret
}

var TestingWrongElectraAggregateAndProofV = func(aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) *electra.AggregateAndProof {
	byts, err := TestingElectraAggregateAndProofV(aggregatorIndex, version).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	ret := &electra.AggregateAndProof{}
	if err := ret.UnmarshalSSZ(byts); err != nil {
		panic(err.Error())
	}
	ret.AggregatorIndex = 100
	return ret
}

var TestingPhase0SignedAggregateAndProof = func(ks *TestKeySet, aggregatorIndex phase0.ValidatorIndex) *phase0.SignedAggregateAndProof {
	agg := TestingPhase0AggregateAndProof(aggregatorIndex)
	return &phase0.SignedAggregateAndProof{
		Message:   agg,
		Signature: signBeaconObject(agg, types.DomainAggregateAndProof, ks),
	}
}

var TestingElectraSignedAggregateAndProofV = func(ks *TestKeySet, aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) *electra.SignedAggregateAndProof {
	agg := TestingElectraAggregateAndProofV(aggregatorIndex, version)
	return &electra.SignedAggregateAndProof{
		Message:   agg,
		Signature: signBeaconObject(agg, types.DomainAggregateAndProof, ks),
	}
}
