package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Attestation Data
// ==================================================

var TestingBlockRoot = phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}

var TestingCommitteeIndex = phase0.CommitteeIndex(3)
var TestingDifferentCommitteeIndex = phase0.CommitteeIndex(4)
var TestingCommitteesAtSlot = uint64(36)
var TestingCommitteeLenght = uint64(128)
var TestingValidatorCommitteeIndex = uint64(11)

var TestingAttestationData = func(version spec.DataVersion) *phase0.AttestationData {
	attData := &phase0.AttestationData{
		Slot:            TestingDutySlotV(version),
		Index:           TestingCommitteeIndex,
		BeaconBlockRoot: TestingBlockRoot,
		Source: &phase0.Checkpoint{
			Epoch: 0,
			Root:  TestingBlockRoot,
		},
		Target: &phase0.Checkpoint{
			Epoch: 1,
			Root:  TestingBlockRoot,
		},
	}

	if version >= spec.DataVersionElectra {
		attData.Index = 0
	}

	return attData
}

var TestingAttestationDataBytes = func(version spec.DataVersion) []byte {
	ret, _ := TestingAttestationData(version).MarshalSSZ()
	return ret
}

var TestingAttestationDataRoot = func(version spec.DataVersion) [32]byte {
	ret, _ := TestingAttestationData(version).HashTreeRoot()
	return ret
}

var TestingAttestationDataForValidatorDuty = func(duty *types.ValidatorDuty) *phase0.AttestationData {
	attData := &phase0.AttestationData{
		Slot:            duty.Slot,
		Index:           duty.CommitteeIndex,
		BeaconBlockRoot: TestBeaconVote.BlockRoot,
		Source:          TestBeaconVote.Source,
		Target:          TestBeaconVote.Target,
	}

	version := VersionBySlot(duty.Slot)
	if version >= spec.DataVersionElectra {
		attData.Index = 0
	}

	return attData
}

var TestingAttestationNextEpochData = func(version spec.DataVersion) *phase0.AttestationData {
	attData := &phase0.AttestationData{
		Slot:            TestingDutySlotNextEpochV(version),
		Index:           TestingCommitteeIndex,
		BeaconBlockRoot: TestingBlockRoot,
		Source: &phase0.Checkpoint{
			Epoch: 0,
			Root:  TestingBlockRoot,
		},
		Target: &phase0.Checkpoint{
			Epoch: 1,
			Root:  TestingBlockRoot,
		},
	}
	if version >= spec.DataVersionElectra {
		attData.Index = 0
	}
	return attData
}
var TestingAttestationNextEpochDataBytes = func(version spec.DataVersion) []byte {
	ret, _ := TestingAttestationNextEpochData(version).MarshalSSZ()
	return ret
}

var TestingWrongAttestationData = func(version spec.DataVersion) *phase0.AttestationData {
	byts, _ := TestingAttestationData(version).MarshalSSZ()
	ret := &phase0.AttestationData{}
	if err := ret.UnmarshalSSZ(byts); err != nil {
		panic(err.Error())
	}
	ret.Slot += 100
	return ret
}

// ==================================================
// Versioned Attestation Response
// ==================================================

var TestingSignedAttestation = func(ks *TestKeySet) *phase0.Attestation {
	duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]
	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
	return &phase0.Attestation{
		Data:            TestingAttestationData(spec.DataVersionPhase0),
		Signature:       signBeaconObject(TestingAttestationData(spec.DataVersionPhase0), types.DomainAttester, ks),
		AggregationBits: aggregationBitfield,
	}
}

var TestingElectraAttestation = func(ks *TestKeySet) *electra.Attestation {
	duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]

	attData := TestingAttestationData(spec.DataVersionElectra)

	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)

	committeeBits := bitfield.NewBitvector64()
	committeeBits.SetBitAt(uint64(duty.CommitteeIndex), true)

	return &electra.Attestation{
		Data:            attData,
		AggregationBits: aggregationBitfield,
		CommitteeBits:   committeeBits,
		Signature:       signBeaconObject(attData, types.DomainAttester, ks),
	}
}

var TestingElectraSingleAttestation = func(ks *TestKeySet) *electra.SingleAttestation {
	duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]

	attData := TestingAttestationData(spec.DataVersionElectra)

	return &electra.SingleAttestation{
		CommitteeIndex: duty.CommitteeIndex,
		AttesterIndex:  TestingValidatorIndex,
		Data:           attData,
		Signature:      signBeaconObject(attData, types.DomainAttester, ks),
	}
}

var TestingVersionedAttestation = func(ks *TestKeySet, version spec.DataVersion) *spec.VersionedAttestation {

	vIdx := phase0.ValidatorIndex(TestingValidatorIndex)

	switch version {
	case spec.DataVersionPhase0:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &vIdx,
			Phase0:         TestingSignedAttestation(ks),
		}

	case spec.DataVersionAltair:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &vIdx,
			Altair:         TestingSignedAttestation(ks),
		}

	case spec.DataVersionBellatrix:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &vIdx,
			Bellatrix:      TestingSignedAttestation(ks),
		}

	case spec.DataVersionCapella:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &vIdx,
			Capella:        TestingSignedAttestation(ks),
		}
	case spec.DataVersionDeneb:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &vIdx,
			Deneb:          TestingSignedAttestation(ks),
		}
	case spec.DataVersionElectra:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &vIdx,
			Electra:        TestingElectraAttestation(ks),
		}
	default:
		panic("unknown data version")
	}
}

var TestingAttestationResponseBeaconObject = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {

	switch version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix, spec.DataVersionCapella, spec.DataVersionDeneb:
		return TestingSignedAttestation(ks)
	case spec.DataVersionElectra:
		return TestingElectraSingleAttestation(ks)
	default:
		panic("unknown data version")
	}
}

// Custom Validator Index

var TestingSignedAttestationForValidatorIndex = func(ks *TestKeySet, validatorIndex phase0.ValidatorIndex) *phase0.Attestation {
	committeeDuty := TestingAttesterDutyForValidator(spec.DataVersionPhase0, validatorIndex)
	duty := committeeDuty.ValidatorDuties[0]
	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
	return &phase0.Attestation{
		Data:            TestingAttestationData(spec.DataVersionPhase0),
		Signature:       signBeaconObject(TestingAttestationData(spec.DataVersionPhase0), types.DomainAttester, ks),
		AggregationBits: aggregationBitfield,
	}
}

var TestingElectraAttestationForValidatorIndex = func(ks *TestKeySet, validatorIndex phase0.ValidatorIndex) *electra.Attestation {
	committeeDuty := TestingAttesterDutyForValidator(spec.DataVersionPhase0, validatorIndex)
	duty := committeeDuty.ValidatorDuties[0]

	attData := TestingAttestationData(spec.DataVersionElectra)

	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)

	committeeBits := bitfield.NewBitvector64()
	committeeBits.SetBitAt(uint64(duty.CommitteeIndex), true)

	return &electra.Attestation{
		Data:            attData,
		AggregationBits: aggregationBitfield,
		CommitteeBits:   committeeBits,
		Signature:       signBeaconObject(attData, types.DomainAttester, ks),
	}
}

var TestingElectraSingleAttestationForValidatorIndex = func(ks *TestKeySet, validatorIndex phase0.ValidatorIndex) *electra.SingleAttestation {

	committeeDuty := TestingAttesterDutyForValidator(spec.DataVersionPhase0, validatorIndex)
	duty := committeeDuty.ValidatorDuties[0]

	attData := TestingAttestationData(spec.DataVersionElectra)

	return &electra.SingleAttestation{
		CommitteeIndex: duty.CommitteeIndex,
		AttesterIndex:  validatorIndex,
		Data:           attData,
		Signature:      signBeaconObject(attData, types.DomainAttester, ks),
	}
}

var TestingSignedAttestationResponseForValidatorIndex = func(ks *TestKeySet, version spec.DataVersion, validatorIndex phase0.ValidatorIndex) *spec.VersionedAttestation {

	switch version {
	case spec.DataVersionPhase0:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &validatorIndex,
			Phase0:         TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}

	case spec.DataVersionAltair:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &validatorIndex,
			Altair:         TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}

	case spec.DataVersionBellatrix:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &validatorIndex,
			Bellatrix:      TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}

	case spec.DataVersionCapella:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &validatorIndex,
			Capella:        TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}
	case spec.DataVersionDeneb:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &validatorIndex,
			Deneb:          TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}
	case spec.DataVersionElectra:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &validatorIndex,
			Electra:        TestingElectraAttestationForValidatorIndex(ks, validatorIndex),
		}
	default:
		panic("unknown data version")
	}
}

var TestingAttestationResponseBeaconObjectForValidatorIndex = func(ks *TestKeySet, version spec.DataVersion, validatorIndex phase0.ValidatorIndex) ssz.HashRoot {

	switch version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix, spec.DataVersionCapella, spec.DataVersionDeneb:
		return TestingSignedAttestationForValidatorIndex(ks, validatorIndex)
	case spec.DataVersionElectra:
		return TestingElectraSingleAttestationForValidatorIndex(ks, validatorIndex)
	default:
		panic("unknown data version")
	}
}

// Custom duty

var TestingSignedAttestationForDuty = func(ks *TestKeySet, duty *types.ValidatorDuty) *phase0.Attestation {
	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
	attData := TestingAttestationDataForValidatorDuty(duty)
	return &phase0.Attestation{
		Data:            attData,
		Signature:       signBeaconObject(attData, types.DomainAttester, ks),
		AggregationBits: aggregationBitfield,
	}
}

var TestingElectraAttestationForDuty = func(ks *TestKeySet, duty *types.ValidatorDuty) *electra.Attestation {
	attData := TestingAttestationDataForValidatorDuty(duty)

	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)

	committeeBits := bitfield.NewBitvector64()
	committeeBits.SetBitAt(uint64(duty.CommitteeIndex), true)

	return &electra.Attestation{
		Data:            attData,
		AggregationBits: aggregationBitfield,
		CommitteeBits:   committeeBits,
		Signature:       signBeaconObject(attData, types.DomainAttester, ks),
	}
}

var TestingElectraSingleAttestationForDuty = func(ks *TestKeySet, duty *types.ValidatorDuty) *electra.SingleAttestation {

	attData := TestingAttestationDataForValidatorDuty(duty)

	return &electra.SingleAttestation{
		CommitteeIndex: duty.CommitteeIndex,
		AttesterIndex:  duty.ValidatorIndex,
		Data:           attData,
		Signature:      signBeaconObject(attData, types.DomainAttester, ks),
	}
}

var TestingSignedAttestationResponseForDuty = func(ks *TestKeySet, version spec.DataVersion, duty *types.ValidatorDuty) *spec.VersionedAttestation {

	switch version {
	case spec.DataVersionPhase0:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &duty.ValidatorIndex,
			Phase0:         TestingSignedAttestationForDuty(ks, duty),
		}

	case spec.DataVersionAltair:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &duty.ValidatorIndex,
			Altair:         TestingSignedAttestationForDuty(ks, duty),
		}

	case spec.DataVersionBellatrix:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &duty.ValidatorIndex,
			Bellatrix:      TestingSignedAttestationForDuty(ks, duty),
		}

	case spec.DataVersionCapella:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &duty.ValidatorIndex,
			Capella:        TestingSignedAttestationForDuty(ks, duty),
		}
	case spec.DataVersionDeneb:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &duty.ValidatorIndex,
			Deneb:          TestingSignedAttestationForDuty(ks, duty),
		}
	case spec.DataVersionElectra:
		return &spec.VersionedAttestation{
			Version:        version,
			ValidatorIndex: &duty.ValidatorIndex,
			Electra:        TestingElectraAttestationForDuty(ks, duty),
		}
	default:
		panic("unknown data version")
	}
}

var TestingAttestationResponseBeaconObjectForDuty = func(ks *TestKeySet, version spec.DataVersion, duty *types.ValidatorDuty) ssz.HashRoot {

	switch version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix, spec.DataVersionCapella, spec.DataVersionDeneb:
		return TestingSignedAttestationForDuty(ks, duty)
	case spec.DataVersionElectra:
		return TestingElectraSingleAttestationForDuty(ks, duty)
	default:
		panic("unknown data version")
	}
}

// SSZ Roots for Key Map

var TestingSignedAttestationSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)
	for _, valKs := range SortedMapKeys(ksMap) {
		ks := valKs.Value
		duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]
		aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
		aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
		ret = append(ret, GetSSZRootNoError(&phase0.Attestation{
			Data:            TestingAttestationData(spec.DataVersionPhase0),
			Signature:       signBeaconObject(TestingAttestationData(spec.DataVersionPhase0), types.DomainAttester, ks),
			AggregationBits: aggregationBitfield,
		}))
	}
	return ret
}

var TestingElectraAttestationSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)
	for _, valKs := range SortedMapKeys(ksMap) {
		ks := valKs.Value
		valIdx := valKs.Key
		committeeDuty := TestingAttesterDutyForValidator(spec.DataVersionElectra, valIdx)
		duty := committeeDuty.ValidatorDuties[0]

		attData := TestingAttestationDataForValidatorDuty(duty)

		aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
		aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)

		committeeBits := bitfield.NewBitvector64()
		committeeBits.SetBitAt(uint64(duty.CommitteeIndex), true)

		att := &electra.Attestation{
			Data:            attData,
			AggregationBits: aggregationBitfield,
			CommitteeBits:   committeeBits,
			Signature:       signBeaconObject(attData, types.DomainAttester, ks),
		}

		ret = append(ret, GetSSZRootNoError(att))
	}
	return ret
}

var TestingElectraSingleAttestationSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)

	for _, valKs := range SortedMapKeys(ksMap) {
		ks := valKs.Value
		valIdx := valKs.Key
		committeeDuty := TestingAttesterDutyForValidator(spec.DataVersionElectra, valIdx)
		duty := committeeDuty.ValidatorDuties[0]

		attData := TestingAttestationDataForValidatorDuty(duty)

		singleAtt := &electra.SingleAttestation{
			CommitteeIndex: duty.CommitteeIndex,
			AttesterIndex:  valIdx,
			Data:           attData,
			Signature:      signBeaconObject(attData, types.DomainAttester, ks),
		}

		ret = append(ret, GetSSZRootNoError(singleAtt))
	}
	return ret
}

var TestingSignedAttestationResponseSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet, version spec.DataVersion) []string {
	switch version {
	case spec.DataVersionPhase0, spec.DataVersionAltair, spec.DataVersionBellatrix, spec.DataVersionCapella, spec.DataVersionDeneb:
		return TestingSignedAttestationSSZRootForKeyMap(ksMap)
	case spec.DataVersionElectra:
		return TestingElectraSingleAttestationSSZRootForKeyMap(ksMap)
	default:
		panic("unknown data version")
	}
}
