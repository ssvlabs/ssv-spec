package testingutils

import (
	"fmt"
	"sort"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

var SupportedAttestationVersions = []spec.DataVersion{spec.DataVersionPhase0, spec.DataVersionElectra}

// ==================================================
// Versioned CommitteeDuty
// ==================================================

func TestingCommitteeDutyForSlot(slot phase0.Slot, attestationValidatorIds []int, syncCommitteeValidatorIds []int) *types.CommitteeDuty {
	return TestingCommitteeDutyWithParams(
		slot,
		attestationValidatorIds,
		syncCommitteeValidatorIds,
		TestingCommitteeIndex,
		TestingCommitteesAtSlot,
		TestingCommitteeLenght,
		TestingValidatorCommitteeIndex)
}

func TestingCommitteeDuty(attestationValidatorIds []int, syncCommitteeValidatorIds []int, version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDutyForSlot(TestingDutySlotV(version), attestationValidatorIds, syncCommitteeValidatorIds)
}
func TestingCommitteeDutyNextEpoch(attestationValidatorIds []int, syncCommitteeValidatorIds []int, version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDutyForSlot(TestingDutySlotNextEpochV(version), attestationValidatorIds, syncCommitteeValidatorIds)
}
func TestingCommitteeDutyInvalid(attestationValidatorIds []int, syncCommitteeValidatorIds []int, version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDutyForSlot(TestingInvalidDutySlotV(version), attestationValidatorIds, syncCommitteeValidatorIds)
}
func TestingCommitteeDutyFirstSlot(attestationValidatorIds []int, syncCommitteeValidatorIds []int) *types.CommitteeDuty {
	return TestingCommitteeDutyForSlot(0, attestationValidatorIds, syncCommitteeValidatorIds)
}

// Create a CommitteeDuty with attestations and sync committee with mixed CommitteeIndexes
func TestingCommitteeDutyWithMixedCommitteeIndexes(attestationValidatorIds []int, syncCommitteeValidatorIds []int, version spec.DataVersion) *types.CommitteeDuty {

	// Sort the validator indexes
	sort.Slice(attestationValidatorIds, func(i, j int) bool {
		return attestationValidatorIds[i] < attestationValidatorIds[j]
	})

	var ret *types.CommitteeDuty
	for i, valIdx := range attestationValidatorIds {
		var duty *types.CommitteeDuty
		// Assign the first half of the validators to a fixed committee index
		if i < len(attestationValidatorIds)/2 {
			duty = TestingCommitteeDuty([]int{valIdx}, nil, version)
		} else {
			// Assign the second half of the validators to a different committee index
			duty = TestingCommitteeDutyWithParams(
				TestingDutySlotV(version),
				[]int{valIdx},
				nil,
				TestingDifferentCommitteeIndex,
				TestingCommitteesAtSlot,
				TestingCommitteeLenght,
				TestingValidatorCommitteeIndex)
		}
		if ret == nil {
			ret = duty
		} else {
			ret.ValidatorDuties = append(ret.ValidatorDuties, duty.ValidatorDuties...)
		}
	}

	sort.Slice(syncCommitteeValidatorIds, func(i, j int) bool {
		return syncCommitteeValidatorIds[i] < syncCommitteeValidatorIds[j]
	})

	for i, valIdx := range syncCommitteeValidatorIds {
		var duty *types.CommitteeDuty
		// Assign the first half of the validators to the a fixed committee index
		if i < len(syncCommitteeValidatorIds)/2 {
			duty = TestingCommitteeDuty(nil, []int{valIdx}, version)
		} else {
			// Assign the second half of the validators to a different committee index
			duty = TestingCommitteeDutyWithParams(
				TestingDutySlotV(version),
				nil,
				[]int{valIdx},
				TestingDifferentCommitteeIndex,
				TestingCommitteesAtSlot,
				TestingCommitteeLenght,
				TestingValidatorCommitteeIndex)
		}
		if ret == nil {
			ret = duty
		} else {
			ret.ValidatorDuties = append(ret.ValidatorDuties, duty.ValidatorDuties...)
		}
	}

	return ret
}

func TestingCommitteeDutyWithParams(
	slot phase0.Slot,
	attestationValidatorIds []int,
	syncCommitteeValidatorIds []int,
	committeeIndex phase0.CommitteeIndex,
	committeesAtSlot uint64,
	committeeLenght uint64,
	validatorCommitteeIndex uint64) *types.CommitteeDuty {

	duties := make([]*types.ValidatorDuty, 0)

	for _, valIdx := range attestationValidatorIds {
		pk := getValPubKeyByValIdx(valIdx)
		duties = append(duties, &types.ValidatorDuty{
			Type:                    types.BNRoleAttester,
			PubKey:                  pk,
			Slot:                    slot,
			ValidatorIndex:          phase0.ValidatorIndex(valIdx),
			CommitteeIndex:          committeeIndex,
			CommitteesAtSlot:        committeesAtSlot,
			CommitteeLength:         committeeLenght,
			ValidatorCommitteeIndex: validatorCommitteeIndex,
		})
	}

	for _, valIdx := range syncCommitteeValidatorIds {
		pk := getValPubKeyByValIdx(valIdx)
		duties = append(duties, &types.ValidatorDuty{
			Type:                          types.BNRoleSyncCommittee,
			PubKey:                        pk,
			Slot:                          slot,
			ValidatorIndex:                phase0.ValidatorIndex(valIdx),
			CommitteeIndex:                committeeIndex,
			CommitteesAtSlot:              committeesAtSlot,
			CommitteeLength:               committeeLenght,
			ValidatorCommitteeIndex:       validatorCommitteeIndex,
			ValidatorSyncCommitteeIndices: TestingContributionProofIndexes,
		})
	}

	return &types.CommitteeDuty{Slot: slot, ValidatorDuties: duties}
}

func getValPubKeyByValIdx(valIdx int) phase0.BLSPubKey {
	return TestingValidatorPubKeyForValidatorIndex(phase0.ValidatorIndex(valIdx))
}

// ==================================================
// Committee duty - Attestation only
// ==================================================

var TestingAttesterDuty = func(version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDuty([]int{TestingValidatorIndex}, nil, version)
}

var TestingAttesterDutyForValidator = func(version spec.DataVersion, validatorIndex phase0.ValidatorIndex) *types.CommitteeDuty {
	return TestingCommitteeDuty([]int{int(validatorIndex)}, nil, version)
}

var TestingAttesterDutyForValidators = func(version spec.DataVersion, validatorIndexLst []int) *types.CommitteeDuty {
	return TestingCommitteeDuty(validatorIndexLst, nil, version)
}

var TestingAttesterDutyNextEpoch = func(version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDutyNextEpoch([]int{TestingValidatorIndex}, nil, version)
}

var TestingAttesterDutyFirstSlot = func() *types.CommitteeDuty {
	return TestingCommitteeDutyFirstSlot([]int{TestingValidatorIndex}, nil)
}

// ==================================================
// Committee duty - Sync Committee only
// ==================================================

var TestingSyncCommitteeDuty = func(version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDuty(nil, []int{TestingValidatorIndex}, version)
}

var TestingSyncCommitteeDutyForValidator = func(version spec.DataVersion, validatorIndex phase0.ValidatorIndex) *types.CommitteeDuty {
	return TestingCommitteeDuty(nil, []int{int(validatorIndex)}, version)
}

var TestingSyncCommitteeDutyForValidators = func(version spec.DataVersion, validatorIndexLst []int) *types.CommitteeDuty {
	return TestingCommitteeDuty(nil, validatorIndexLst, version)
}

var TestingSyncCommitteeDutyNextEpoch = func(version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDutyNextEpoch(nil, []int{TestingValidatorIndex}, version)
}

var TestingSyncCommitteeDutyFirstSlot = func() *types.CommitteeDuty {
	return TestingCommitteeDutyFirstSlot(nil, []int{TestingValidatorIndex})
}

// ==================================================
// Committee duty - Attestation and Sync Committee
// ==================================================

var TestingAttesterAndSyncCommitteeDuties = func(version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDuty([]int{TestingValidatorIndex}, []int{TestingValidatorIndex}, version)
}
var TestingAttesterAndSyncCommitteeDutiesNextEpoch = func(version spec.DataVersion) *types.CommitteeDuty {
	return TestingCommitteeDutyNextEpoch([]int{TestingValidatorIndex}, []int{TestingValidatorIndex}, version)
}
var TestingAttesterAndSyncCommitteeDutiesFirstSlot = func() *types.CommitteeDuty {
	return TestingCommitteeDutyFirstSlot([]int{TestingValidatorIndex}, []int{TestingValidatorIndex})
}

// ==================================================
// Beacon Roots for Committee duty
// ==================================================

var TestingSignedCommitteeBeaconObjectSSZRoot = func(duty *types.CommitteeDuty, ksMap map[phase0.ValidatorIndex]*TestKeySet, version spec.DataVersion) []string {
	ret := make([]string, 0)
	for _, validatorDuty := range duty.ValidatorDuties {

		ks := ksMap[validatorDuty.ValidatorIndex]

		if validatorDuty.Type == types.BNRoleAttester {

			attResponse := TestingAttestationResponseBeaconObjectForDuty(ks, version, validatorDuty)
			ret = append(ret, GetSSZRootNoError(attResponse))
		} else if validatorDuty.Type == types.BNRoleSyncCommittee {
			ret = append(ret, GetSSZRootNoError(&altair.SyncCommitteeMessage{
				Slot:            validatorDuty.Slot,
				BeaconBlockRoot: TestingBlockRoot,
				ValidatorIndex:  validatorDuty.ValidatorIndex,
				Signature:       signBeaconObject(types.SSZBytes(TestingBlockRoot[:]), types.DomainSyncCommittee, ks),
			}))
		} else {
			panic(fmt.Sprintf("type %v not expected", validatorDuty.Type))
		}
	}
	return ret
}
