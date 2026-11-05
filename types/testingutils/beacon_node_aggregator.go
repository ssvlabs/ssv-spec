package testingutils

import (
	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec"
	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

var SupportedAggregatorVersions = []spec.DataVersion{gloas.DataVersionGloas}

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
	d := TestingAggregatorCommitteeDutyOnlyAggregator(gloas.DataVersionGloas)
	d.Slot = 0
	for i := range d.ValidatorDuties {
		d.ValidatorDuties[i].Slot = 0
	}
	return d
}

// ==================================================
// Versioned AggregateAndProof
// ==================================================
//
// The Gloas AggregateAndProof (SIP #94 §2) serializes like the Electra one but has a different hash tree root. The
// attestation data carries the payload-status index and the duty slot, so fixtures must match the data the runner
// aggregates over.

var TestingAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	return TestingGloasAggregateAndProofV(aggregatorIndex, version)
}

var TestingVersionedSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) *spec.VersionedSignedAggregateAndProof {
	return &spec.VersionedSignedAggregateAndProof{
		Version: version,
		Gloas:   TestingGloasSignedAggregateAndProofV(ks, TestingValidatorIndex, version),
	}
}

var TestingSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) types.HashRoot {
	return TestingGloasSignedAggregateAndProofV(ks, TestingValidatorIndex, version)
}

var TestingAggregateAndProofBytesV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) []byte {
	return TestingGloasAggregateAndProofBytesV(aggregatorIndex, version)
}

var TestingWrongAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	return TestingWrongGloasAggregateAndProofV(aggregatorIndex, version)
}

var TestingGloasAggregateAndProofV = func(aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) *eth2gloas.AggregateAndProof {
	return &eth2gloas.AggregateAndProof{
		AggregatorIndex: aggregatorIndex,
		SelectionProof:  phase0.BLSSignature{},
		Aggregate: &eth2gloas.Attestation{
			AggregationBits: bitfield.NewBitlist(128),
			Signature:       phase0.BLSSignature{},
			Data:            TestingAttestationData(version),
			CommitteeBits:   bitfield.NewBitvector64(),
		},
	}
}

var TestingGloasAggregateAndProofBytesV = func(aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) []byte {
	ret, _ := TestingGloasAggregateAndProofV(aggregatorIndex, version).MarshalSSZ()
	return ret
}

var TestingWrongGloasAggregateAndProofV = func(aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) *eth2gloas.AggregateAndProof {
	byts, err := TestingGloasAggregateAndProofV(aggregatorIndex, version).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	ret := &eth2gloas.AggregateAndProof{}
	if err := ret.UnmarshalSSZ(byts); err != nil {
		panic(err.Error())
	}
	ret.AggregatorIndex = 100
	return ret
}

var TestingGloasSignedAggregateAndProofV = func(ks *TestKeySet, aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) *eth2gloas.SignedAggregateAndProof {
	agg := TestingGloasAggregateAndProofV(aggregatorIndex, version)
	return &eth2gloas.SignedAggregateAndProof{
		Message:   agg,
		Signature: signBeaconObject(agg, types.DomainAggregateAndProof, ks),
	}
}
