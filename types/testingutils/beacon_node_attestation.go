package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
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

var TestingAttestationData = &phase0.AttestationData{
	Slot:            TestingDutySlot,
	Index:           0, // EIP-7549: Index should be set to 0
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

var TestingAttestationDataRoot, _ = TestingAttestationData.HashTreeRoot()

var TestingAttestationDataForValidatorDuty = func(duty *types.ValidatorDuty) *phase0.AttestationData {
	return &phase0.AttestationData{
		Slot:            duty.Slot,
		Index:           0, // EIP-7549: Index should be set to 0
		BeaconBlockRoot: TestBeaconVote.BlockRoot,
		Source:          TestBeaconVote.Source,
		Target:          TestBeaconVote.Target,
	}
}

var TestingAttestationDataBytes = func() []byte {
	ret, _ := TestingAttestationData.MarshalSSZ()
	return ret
}()

var TestingAttestationNextEpochData = &phase0.AttestationData{
	Slot:            TestingDutySlot2,
	Index:           0, // EIP-7549: Index should be set to 0
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
var TestingAttestationNextEpochDataBytes = func() []byte {
	ret, _ := TestingAttestationNextEpochData.MarshalSSZ()
	return ret
}()

var TestingWrongAttestationData = func() *phase0.AttestationData {
	byts, _ := TestingAttestationData.MarshalSSZ()
	ret := &phase0.AttestationData{}
	if err := ret.UnmarshalSSZ(byts); err != nil {
		panic(err.Error())
	}
	ret.Slot = 100
	return ret
}()

// ==================================================
// Versioned Attestation Response
// ==================================================

var TestingSignedAttestation = func(ks *TestKeySet) *phase0.Attestation {
	duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]
	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
	return &phase0.Attestation{
		Data:            TestingAttestationData,
		Signature:       signBeaconObject(TestingAttestationData, types.DomainAttester, ks),
		AggregationBits: aggregationBitfield,
	}
}

var TestingElectraSingleAttestation = func(ks *TestKeySet) *electra.SingleAttestation {
	duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]
	return &electra.SingleAttestation{
		CommitteeIndex: duty.CommitteeIndex,
		AttesterIndex:  duty.ValidatorIndex,
		Data:           TestingAttestationData,
		Signature:      signBeaconObject(TestingAttestationData, types.DomainAttester, ks),
	}
}

var TestingSignedAttestationResponse = func(ks *TestKeySet, version spec.DataVersion) *types.VersionedAttestationResponse {

	switch version {
	case spec.DataVersionPhase0:
		return &types.VersionedAttestationResponse{
			Version: version,
			Phase0:  TestingSignedAttestation(ks),
		}

	case spec.DataVersionAltair:
		return &types.VersionedAttestationResponse{
			Version: version,
			Altair:  TestingSignedAttestation(ks),
		}

	case spec.DataVersionBellatrix:
		return &types.VersionedAttestationResponse{
			Version:   version,
			Bellatrix: TestingSignedAttestation(ks),
		}

	case spec.DataVersionCapella:
		return &types.VersionedAttestationResponse{
			Version: version,
			Capella: TestingSignedAttestation(ks),
		}
	case spec.DataVersionDeneb:
		return &types.VersionedAttestationResponse{
			Version: version,
			Deneb:   TestingSignedAttestation(ks),
		}
	case spec.DataVersionElectra:
		return &types.VersionedAttestationResponse{
			Version: version,
			Electra: TestingElectraSingleAttestation(ks),
		}
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
		Data:            TestingAttestationData,
		Signature:       signBeaconObject(TestingAttestationData, types.DomainAttester, ks),
		AggregationBits: aggregationBitfield,
	}
}

var TestingElectraSingleAttestationForValidatorIndex = func(ks *TestKeySet, validatorIndex phase0.ValidatorIndex) *electra.SingleAttestation {
	committeeDuty := TestingAttesterDutyForValidator(spec.DataVersionPhase0, validatorIndex)
	duty := committeeDuty.ValidatorDuties[0]
	return &electra.SingleAttestation{
		CommitteeIndex: duty.CommitteeIndex,
		AttesterIndex:  validatorIndex,
		Data:           TestingAttestationData,
		Signature:      signBeaconObject(TestingAttestationData, types.DomainAttester, ks),
	}
}

var TestingSignedAttestationResponseForValidatorIndex = func(ks *TestKeySet, version spec.DataVersion, validatorIndex phase0.ValidatorIndex) *types.VersionedAttestationResponse {

	switch version {
	case spec.DataVersionPhase0:
		return &types.VersionedAttestationResponse{
			Version: version,
			Phase0:  TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}

	case spec.DataVersionAltair:
		return &types.VersionedAttestationResponse{
			Version: version,
			Altair:  TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}

	case spec.DataVersionBellatrix:
		return &types.VersionedAttestationResponse{
			Version:   version,
			Bellatrix: TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}

	case spec.DataVersionCapella:
		return &types.VersionedAttestationResponse{
			Version: version,
			Capella: TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}
	case spec.DataVersionDeneb:
		return &types.VersionedAttestationResponse{
			Version: version,
			Deneb:   TestingSignedAttestationForValidatorIndex(ks, validatorIndex),
		}
	case spec.DataVersionElectra:
		return &types.VersionedAttestationResponse{
			Version: version,
			Electra: TestingElectraSingleAttestationForValidatorIndex(ks, validatorIndex),
		}
	default:
		panic("unknown data version")
	}
}

// SSZ Roots for Key Map

var TestingSignedAttestationSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)
	for _, ks := range ksMap {
		duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]
		aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
		aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
		ret = append(ret, GetSSZRootNoError(&phase0.Attestation{
			Data:            TestingAttestationData,
			Signature:       signBeaconObject(TestingAttestationData, types.DomainAttester, ks),
			AggregationBits: aggregationBitfield,
		}))
	}
	return ret
}

var TestingElectraSingleAttestationSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)
	for valIdx, ks := range ksMap {
		duty := TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0]
		ret = append(ret, GetSSZRootNoError(&electra.SingleAttestation{
			CommitteeIndex: duty.CommitteeIndex,
			AttesterIndex:  valIdx,
			Data:           TestingAttestationData,
			Signature:      signBeaconObject(TestingAttestationData, types.DomainAttester, ks),
		}))
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
