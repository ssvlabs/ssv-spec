package testingutils

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Sync Committee Message
// ==================================================

var TestingSyncCommitteeBlockRoot = TestingBlockRoot
var TestingSyncCommitteeWrongBlockRoot = phase0.Root{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

var TestingSignedSyncCommitteeBlockRoot = func(ks *TestKeySet, version spec.DataVersion) *altair.SyncCommitteeMessage {
	return &altair.SyncCommitteeMessage{
		Slot:            TestingDutySlotV(version),
		BeaconBlockRoot: TestingSyncCommitteeBlockRoot,
		ValidatorIndex:  TestingValidatorIndex,
		Signature:       signBeaconObject(types.SSZBytes(TestingSyncCommitteeBlockRoot[:]), types.DomainSyncCommittee, ks),
	}
}

var TestingSignedSyncCommitteeBlockRootForValidatorIndex = func(ks *TestKeySet, validatorIndex phase0.ValidatorIndex, version spec.DataVersion) *altair.SyncCommitteeMessage {
	return &altair.SyncCommitteeMessage{
		Slot:            TestingDutySlotV(version),
		BeaconBlockRoot: TestingSyncCommitteeBlockRoot,
		ValidatorIndex:  validatorIndex,
		Signature:       signBeaconObject(types.SSZBytes(TestingSyncCommitteeBlockRoot[:]), types.DomainSyncCommittee, ks),
	}
}

var TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet, version spec.DataVersion) []string {
	ret := make([]string, 0)
	for _, valKs := range SortedMapKeys(ksMap) {
		ks := valKs.Value
		valIdx := valKs.Key
		ret = append(ret, GetSSZRootNoError(&altair.SyncCommitteeMessage{
			Slot:            TestingDutySlotV(version),
			BeaconBlockRoot: TestingBlockRoot,
			ValidatorIndex:  valIdx,
			Signature:       signBeaconObject(types.SSZBytes(TestingBlockRoot[:]), types.DomainSyncCommittee, ks),
		}))
	}
	return ret
}

// ==================================================
// Sync Committee Contribution Objects
// ==================================================

var TestingContributionProofIndexes = []uint64{0, 1, 2}
var TestingContributionProofsSigned = func() []phase0.BLSSignature {
	// signed with 3515c7d08e5affd729e9579f7588d30f2342ee6f6a9334acf006345262162c6f
	byts1, _ := hex.DecodeString("b18833bb7549ec33e8ac414ba002fd45bb094ca300bd24596f04a434a89beea462401da7c6b92fb3991bd17163eb603604a40e8dd6781266c990023446776ff42a9313df26a0a34184a590e57fa4003d610c2fa214db4e7dec468592010298bc")
	byts2, _ := hex.DecodeString("9094342c95146554df849dc20f7425fca692dacee7cb45258ddd264a8e5929861469fda3d1567b9521cba83188ffd61a0dbe6d7180c7a96f5810d18db305e9143772b766d368aa96d3751f98d0ce2db9f9e6f26325702088d87f0de500c67c68")
	byts3, _ := hex.DecodeString("a7f88ce43eff3aa8cdd2e3957c5bead4e21353fbecac6079a5398d03019bc45ff7c951785172deee70e9bc5abbc8ca6a0f0441e9d4cc9da74c31121357f7d7c7de9533f6f457da493e3314e22d554ab76613e469b050e246aff539a33807197c")

	ret := make([]phase0.BLSSignature, 0)
	for _, byts := range [][]byte{byts1, byts2, byts3} {
		b := phase0.BLSSignature{}
		copy(b[:], byts)
		ret = append(ret, b)
	}
	return ret
}()

var TestingSyncCommitteeContributions = []*altair.SyncCommitteeContribution{
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 0,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         phase0.BLSSignature{},
	},
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 1,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         phase0.BLSSignature{},
	},
	{
		Slot:              TestingDutySlot,
		BeaconBlockRoot:   TestingSyncCommitteeBlockRoot,
		SubcommitteeIndex: 2,
		AggregationBits:   bitfield.NewBitvector128(),
		Signature:         phase0.BLSSignature{},
	},
}

var TestingContributionsData = func() types.Contributions {
	d := types.Contributions{}
	d = append(d, &types.Contribution{
		SelectionProofSig: TestingContributionProofsSigned[0],
		Contribution:      *TestingSyncCommitteeContributions[0],
	})
	d = append(d, &types.Contribution{
		SelectionProofSig: TestingContributionProofsSigned[1],
		Contribution:      *TestingSyncCommitteeContributions[1],
	})
	d = append(d, &types.Contribution{
		SelectionProofSig: TestingContributionProofsSigned[2],
		Contribution:      *TestingSyncCommitteeContributions[2],
	})
	return d
}()

var TestingContributionsDataBytes = func() []byte {
	ret, _ := TestingContributionsData.MarshalSSZ()
	return ret
}()

var TestingSignedSyncCommitteeContributions = func(
	contrib *altair.SyncCommitteeContribution,
	proof phase0.BLSSignature,
	ks *TestKeySet) *altair.SignedContributionAndProof {
	msg := &altair.ContributionAndProof{
		AggregatorIndex: TestingValidatorIndex,
		Contribution:    contrib,
		SelectionProof:  proof,
	}
	return &altair.SignedContributionAndProof{
		Message:   msg,
		Signature: signBeaconObject(msg, types.DomainContributionAndProof, ks),
	}
}

// ==================================================
// Sync Committee Contribution Duty
// ==================================================

var TestingSyncCommitteeContributionDuty = types.ValidatorDuty{
	Type:                          types.BNRoleSyncCommitteeContribution,
	PubKey:                        TestingValidatorPubKey,
	Slot:                          TestingDutySlot,
	ValidatorIndex:                TestingValidatorIndex,
	CommitteeIndex:                3,
	CommitteesAtSlot:              36,
	CommitteeLength:               128,
	ValidatorCommitteeIndex:       11,
	ValidatorSyncCommitteeIndices: TestingContributionProofIndexes,
}

var TestingSyncCommitteeContributionDutyFirstSlot = types.ValidatorDuty{
	Type:                          types.BNRoleSyncCommitteeContribution,
	PubKey:                        TestingValidatorPubKey,
	Slot:                          0,
	ValidatorIndex:                TestingValidatorIndex,
	CommitteeIndex:                3,
	CommitteesAtSlot:              36,
	CommitteeLength:               128,
	ValidatorCommitteeIndex:       11,
	ValidatorSyncCommitteeIndices: TestingContributionProofIndexes,
}

var TestingSyncCommitteeContributionNexEpochDuty = types.ValidatorDuty{
	Type:                          types.BNRoleSyncCommitteeContribution,
	PubKey:                        TestingValidatorPubKey,
	Slot:                          TestingDutySlot2,
	ValidatorIndex:                TestingValidatorIndex,
	CommitteeIndex:                3,
	CommitteesAtSlot:              36,
	CommitteeLength:               128,
	ValidatorCommitteeIndex:       11,
	ValidatorSyncCommitteeIndices: TestingContributionProofIndexes,
}
