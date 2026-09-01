package testingutils

import (
	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
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
// Gloas reuses the Electra aggregate shape (SIP #94 §2); the attestation data carries the payload-status
// index and the duty slot, so fixtures must match the data the runner aggregates over.

var TestingAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	return TestingElectraAggregateAndProofV(aggregatorIndex, version)
}

var TestingVersionedSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) *spec.VersionedSignedAggregateAndProof {
	return &spec.VersionedSignedAggregateAndProof{
		Version: version,
		Electra: TestingElectraSignedAggregateAndProofV(ks, TestingValidatorIndex, version),
	}
}

var TestingSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) types.HashRoot {
	return TestingElectraSignedAggregateAndProofV(ks, TestingValidatorIndex, version)
}

var TestingAggregateAndProofBytesV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) []byte {
	return TestingElectraAggregateAndProofBytesV(aggregatorIndex, version)
}

var TestingWrongAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	return TestingWrongElectraAggregateAndProofV(aggregatorIndex, version)
}

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

var TestingElectraSignedAggregateAndProofV = func(ks *TestKeySet, aggregatorIndex phase0.ValidatorIndex, version spec.DataVersion) *electra.SignedAggregateAndProof {
	agg := TestingElectraAggregateAndProofV(aggregatorIndex, version)
	return &electra.SignedAggregateAndProof{
		Message:   agg,
		Signature: signBeaconObject(agg, types.DomainAggregateAndProof, ks),
	}
}
