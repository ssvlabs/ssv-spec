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

var TestingAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	if version == spec.DataVersionElectra {
		return TestingElectraAggregateAndProof(aggregatorIndex)
	} else {
		return TestingPhase0AggregateAndProof(aggregatorIndex)
	}
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
			Electra: TestingElectraSignedAggregateAndProof(ks, TestingValidatorIndex),
		}
	default:
		panic("unknown data version")
	}
}

var TestingSignedAggregateAndProof = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	switch version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix, spec.DataVersionCapella, spec.DataVersionDeneb:
		return TestingPhase0SignedAggregateAndProof(ks, TestingValidatorIndex)
	case spec.DataVersionElectra:
		return TestingElectraSignedAggregateAndProof(ks, TestingValidatorIndex)
	default:
		panic("unknown data version")
	}
}

var TestingAggregateAndProofBytesV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) []byte {
	if version == spec.DataVersionElectra {
		return TestingElectraAggregateAndProofBytes(aggregatorIndex)
	} else {
		return TestingPhase0AggregateAndProofBytes(aggregatorIndex)
	}
}

var TestingWrongAggregateAndProofV = func(version spec.DataVersion, aggregatorIndex phase0.ValidatorIndex) ssz.Marshaler {
	if version == spec.DataVersionElectra {
		return TestingWrongElectraAggregateAndProof(aggregatorIndex)
	} else {
		return TestingWrongPhase0AggregateAndProof(aggregatorIndex)
	}
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

var TestingElectraAggregateAndProof = func(aggregatorIndex phase0.ValidatorIndex) *electra.AggregateAndProof {
	return &electra.AggregateAndProof{
		AggregatorIndex: aggregatorIndex,
		SelectionProof:  phase0.BLSSignature{},
		Aggregate: &electra.Attestation{
			AggregationBits: bitfield.NewBitlist(128),
			Signature:       phase0.BLSSignature{},
			Data:            TestingAttestationData(spec.DataVersionElectra),
			CommitteeBits:   bitfield.NewBitvector64(),
		},
	}
}
var TestingElectraAggregateAndProofBytes = func(aggregatorIndex phase0.ValidatorIndex) []byte {
	ret, _ := TestingElectraAggregateAndProof(aggregatorIndex).MarshalSSZ()
	return ret
}

var TestingWrongElectraAggregateAndProof = func(aggregatorIndex phase0.ValidatorIndex) *electra.AggregateAndProof {
	byts, err := TestingElectraAggregateAndProof(aggregatorIndex).MarshalSSZ()
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

var TestingElectraSignedAggregateAndProof = func(ks *TestKeySet, aggregatorIndex phase0.ValidatorIndex) *electra.SignedAggregateAndProof {
	agg := TestingElectraAggregateAndProof(aggregatorIndex)
	return &electra.SignedAggregateAndProof{
		Message:   agg,
		Signature: signBeaconObject(agg, types.DomainAggregateAndProof, ks),
	}
}
