package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

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

// TestingAttestationData carries the Gloas payload-status index (SIP #94 §2: 1 = payload present, matching
// TestBeaconVote).
var TestingAttestationData = func(version spec.DataVersion) *phase0.AttestationData {
	return &phase0.AttestationData{
		Slot:            TestingDutySlotV(version),
		Index:           TestBeaconVote.AttestationDataIndex,
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
}

var TestingAttestationDataBytes = func(version spec.DataVersion) []byte {
	ret, _ := TestingAttestationData(version).MarshalSSZ()
	return ret
}

var TestingAttestationDataForValidatorDuty = func(duty *types.ValidatorDuty) *phase0.AttestationData {
	return &phase0.AttestationData{
		Slot:            duty.Slot,
		Index:           TestBeaconVote.AttestationDataIndex, // SIP #94 §2: the payload-status index
		BeaconBlockRoot: TestBeaconVote.BlockRoot,
		Source:          TestBeaconVote.Source,
		Target:          TestBeaconVote.Target,
	}
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
//
// The runner submits the Electra-shaped attestation (Gloas reuses it, SIP #94 §2), which the testing beacon
// node records as the on-wire SingleAttestation.

var TestingAttestationResponseBeaconObject = func(ks *TestKeySet, version spec.DataVersion) ssz.HashRoot {
	return TestingAttestationResponseBeaconObjectForValidatorIndex(ks, version, TestingValidatorIndex)
}

var TestingAttestationResponseBeaconObjectForValidatorIndex = func(ks *TestKeySet, version spec.DataVersion, validatorIndex phase0.ValidatorIndex) ssz.HashRoot {
	duty := TestingAttesterDutyForValidator(version, validatorIndex).ValidatorDuties[0]
	attData := TestingAttestationData(version)
	return &electra.SingleAttestation{
		CommitteeIndex: duty.CommitteeIndex,
		AttesterIndex:  validatorIndex,
		Data:           attData,
		Signature:      signBeaconObject(attData, types.DomainAttester, ks),
	}
}

var TestingAttestationResponseBeaconObjectForDuty = func(ks *TestKeySet, version spec.DataVersion, duty *types.ValidatorDuty) ssz.HashRoot {
	attData := TestingAttestationDataForValidatorDuty(duty)
	return &electra.SingleAttestation{
		CommitteeIndex: duty.CommitteeIndex,
		AttesterIndex:  duty.ValidatorIndex,
		Data:           attData,
		Signature:      signBeaconObject(attData, types.DomainAttester, ks),
	}
}

// TestingSignedAttestationResponseSSZRootForKeyMap builds the expected submitted-attestation roots for the
// validators in ksMap.
var TestingSignedAttestationResponseSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet, version spec.DataVersion) []string {
	ret := make([]string, 0)
	for _, valKs := range SortedMapKeys(ksMap) {
		duty := TestingAttesterDutyForValidator(version, valKs.Key).ValidatorDuties[0]
		ret = append(ret, GetSSZRootNoError(TestingAttestationResponseBeaconObjectForDuty(valKs.Value, version, duty)))
	}
	return ret
}
