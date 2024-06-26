package testingutils

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/types"
)

var signBeaconObject = func(obj ssz.HashRoot, domainType phase0.DomainType, ks *TestKeySet) phase0.BLSSignature {
	domain, _ := NewTestingBeaconNode().DomainData(1, domainType)
	ret, _, _ := NewTestingKeyManager().SignBeaconObject(obj, domain, ks.ValidatorPK.Serialize(), domainType)

	blsSig := phase0.BLSSignature{}
	copy(blsSig[:], ret)

	return blsSig
}

func GetSSZRootNoError(obj ssz.HashRoot) string {
	r, _ := obj.HashTreeRoot()
	return hex.EncodeToString(r[:])
}

var TestingBlockRoot = phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}

var TestingCommitteeIndex = phase0.CommitteeIndex(3)
var TestingDifferentCommitteeIndex = phase0.CommitteeIndex(4)
var TestingCommitteesAtSlot = uint64(36)
var TestingCommitteeLenght = uint64(128)
var TestingValidatorCommitteeIndex = uint64(11)

// Beacon vote

var TestBeaconVote = types.BeaconVote{
	BlockRoot: TestingBlockRoot,
	Source: &phase0.Checkpoint{
		Epoch: 0,
		Root:  TestingBlockRoot,
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  TestingBlockRoot,
	},
}

var TestBeaconVoteNextEpoch = types.BeaconVote{
	BlockRoot: TestingBlockRoot,
	Source: &phase0.Checkpoint{
		Epoch: 0,
		Root:  TestingBlockRoot,
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  TestingBlockRoot,
	},
}

var TestWrongBeaconVote = types.BeaconVote{
	BlockRoot: phase0.Root{1, 2, 3, 4},
	Source: &phase0.Checkpoint{
		Epoch: 2,
		Root:  phase0.Root{1, 2, 3, 4},
	},
	Target: &phase0.Checkpoint{
		Epoch: 1,
		Root:  phase0.Root{1, 2, 3, 5},
	},
}

var TestingAttestationData = &phase0.AttestationData{
	Slot:            TestingDutySlot,
	Index:           3,
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
		Index:           duty.CommitteeIndex,
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
	Index:           3,
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

var TestingSignedAttestation = func(ks *TestKeySet) *phase0.Attestation {
	duty := TestingAttesterDuty.ValidatorDuties[0]
	aggregationBitfield := bitfield.NewBitlist(duty.CommitteeLength)
	aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
	return &phase0.Attestation{
		Data:            TestingAttestationData,
		Signature:       signBeaconObject(TestingAttestationData, types.DomainAttester, ks),
		AggregationBits: aggregationBitfield,
	}
}

var TestingSignedAttestationSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)
	for _, ks := range ksMap {
		duty := TestingAttesterDuty.ValidatorDuties[0]
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

var TestingSignedCommitteeBeaconObjectSSZRoot = func(duty *types.CommitteeDuty, ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)
	for _, validatorDuty := range duty.ValidatorDuties {

		ks := ksMap[validatorDuty.ValidatorIndex]

		if validatorDuty.Type == types.BNRoleAttester {
			attData := TestingAttestationDataForValidatorDuty(validatorDuty)
			aggregationBitfield := bitfield.NewBitlist(validatorDuty.CommitteeLength)
			aggregationBitfield.SetBitAt(validatorDuty.ValidatorCommitteeIndex, true)
			ret = append(ret, GetSSZRootNoError(&phase0.Attestation{
				Data:            attData,
				Signature:       signBeaconObject(attData, types.DomainAttester, ks),
				AggregationBits: aggregationBitfield,
			}))
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

var TestingAggregateAndProof = &phase0.AggregateAndProof{
	AggregatorIndex: 1,
	SelectionProof:  phase0.BLSSignature{},
	Aggregate: &phase0.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Signature:       phase0.BLSSignature{},
		Data:            TestingAttestationData,
	},
}
var TestingAggregateAndProofBytes = func() []byte {
	ret, _ := TestingAggregateAndProof.MarshalSSZ()
	return ret
}()

var TestingWrongAggregateAndProof = func() *phase0.AggregateAndProof {
	byts, err := TestingAggregateAndProof.MarshalSSZ()
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

var TestingSignedAggregateAndProof = func(ks *TestKeySet) *phase0.SignedAggregateAndProof {
	return &phase0.SignedAggregateAndProof{
		Message:   TestingAggregateAndProof,
		Signature: signBeaconObject(TestingAggregateAndProof, types.DomainAggregateAndProof, ks),
	}
}

const (
	TestingDutySlot            = 12
	TestingDutySlot2           = 50
	TestingDutyEpoch           = 0
	TestingValidatorIndex      = 1
	TestingWrongValidatorIndex = 100

	UnknownDutyType = 100
)

var TestingSyncCommitteeBlockRoot = TestingBlockRoot
var TestingSyncCommitteeWrongBlockRoot = phase0.Root{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
var TestingSignedSyncCommitteeBlockRoot = func(ks *TestKeySet) *altair.SyncCommitteeMessage {
	return &altair.SyncCommitteeMessage{
		Slot:            TestingDutySlot,
		BeaconBlockRoot: TestingSyncCommitteeBlockRoot,
		ValidatorIndex:  TestingValidatorIndex,
		Signature:       signBeaconObject(types.SSZBytes(TestingSyncCommitteeBlockRoot[:]), types.DomainSyncCommittee, ks),
	}
}

var TestingSignedSyncCommitteeBlockRootSSZRootForKeyMap = func(ksMap map[phase0.ValidatorIndex]*TestKeySet) []string {
	ret := make([]string, 0)
	for valIdx, ks := range ksMap {
		ret = append(ret, GetSSZRootNoError(&altair.SyncCommitteeMessage{
			Slot:            TestingDutySlot,
			BeaconBlockRoot: TestingBlockRoot,
			ValidatorIndex:  valIdx,
			Signature:       signBeaconObject(types.SSZBytes(TestingBlockRoot[:]), types.DomainSyncCommittee, ks),
		}))
	}
	return ret
}

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

var TestingFeeRecipient = bellatrix.ExecutionAddress(ethAddressFromHex("535953b5a6040074948cf185eaa7d2abbd66808f"))
var TestingValidatorRegistration = &v1.ValidatorRegistration{
	FeeRecipient: TestingFeeRecipient,
	GasLimit:     types.DefaultGasLimit,
	Timestamp:    types.PraterNetwork.EpochStartTime(TestingDutyEpoch),
	Pubkey:       TestingValidatorPubKey,
}
var TestingValidatorRegistrationWrong = &v1.ValidatorRegistration{
	FeeRecipient: TestingFeeRecipient,
	GasLimit:     5,
	Timestamp:    types.PraterNetwork.EpochStartTime(TestingDutyEpoch),
	Pubkey:       TestingValidatorPubKey,
}

// TestingValidatorRegistrationBySlot receives a slot and calculates the correct timestamp
func TestingValidatorRegistrationBySlot(slot phase0.Slot) *v1.ValidatorRegistration {
	epoch := types.PraterNetwork.EstimatedEpochAtSlot(slot)
	return &v1.ValidatorRegistration{
		FeeRecipient: TestingFeeRecipient,
		GasLimit:     types.DefaultGasLimit,
		Timestamp:    types.PraterNetwork.EpochStartTime(epoch),
		Pubkey:       TestingValidatorPubKey,
	}
}

var TestingVoluntaryExit = &phase0.VoluntaryExit{
	Epoch:          0,
	ValidatorIndex: TestingValidatorIndex,
}
var TestingVoluntaryExitWrong = &phase0.VoluntaryExit{
	Epoch:          1,
	ValidatorIndex: TestingValidatorIndex,
}
var TestingSignedVoluntaryExit = func(ks *TestKeySet) *phase0.SignedVoluntaryExit {
	return &phase0.SignedVoluntaryExit{
		Message:   TestingVoluntaryExit,
		Signature: signBeaconObject(TestingVoluntaryExit, types.DomainVoluntaryExit, ks),
	}
}

// TestingVoluntaryExitBySlot receives a slot and calculates the correct epoch
func TestingVoluntaryExitBySlot(slot phase0.Slot) *phase0.VoluntaryExit {
	epoch := types.PraterNetwork.EstimatedEpochAtSlot(slot)
	return &phase0.VoluntaryExit{
		Epoch:          epoch,
		ValidatorIndex: TestingValidatorIndex,
	}
}

// TestingProposerDutyFirstSlot
var TestingProposerDutyFirstSlot = types.ValidatorDuty{
	Type:           types.BNRoleProposer,
	PubKey:         TestingValidatorPubKey,
	Slot:           0,
	ValidatorIndex: TestingValidatorIndex,
}

func getValPubKeyByValIdx(valIdx int) phase0.BLSPubKey {
	return TestingValidatorPubKeyForValidatorIndex(phase0.ValidatorIndex(valIdx))
}

func TestingCommitteeDuty(slot phase0.Slot, attestationValidatorIds []int, syncCommitteeValidatorIds []int) *types.CommitteeDuty {
	return TestingCommitteeDutyWithParams(slot, attestationValidatorIds, syncCommitteeValidatorIds, TestingCommitteeIndex, TestingCommitteesAtSlot, TestingCommitteeLenght, TestingValidatorCommitteeIndex)
}

func TestingCommitteeDutyWithMixedCommitteeIndexes(slot phase0.Slot, attestationValidatorIds []int, syncCommitteeValidatorIds []int) *types.CommitteeDuty {
	sort.Slice(attestationValidatorIds, func(i, j int) bool {
		return attestationValidatorIds[i] < attestationValidatorIds[j]
	})

	var ret *types.CommitteeDuty
	for i, valIdx := range attestationValidatorIds {
		var duty *types.CommitteeDuty
		if i < len(attestationValidatorIds)/2 {
			duty = TestingCommitteeDuty(slot, []int{valIdx}, nil)
		} else {
			duty = TestingCommitteeDutyWithParams(slot, []int{valIdx}, nil, TestingDifferentCommitteeIndex, TestingCommitteesAtSlot, TestingCommitteeLenght, TestingValidatorCommitteeIndex)
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
		if i < len(syncCommitteeValidatorIds)/2 {
			duty = TestingCommitteeDuty(slot, nil, []int{valIdx})
		} else {
			duty = TestingCommitteeDutyWithParams(slot, nil, []int{valIdx}, TestingDifferentCommitteeIndex, TestingCommitteesAtSlot, TestingCommitteeLenght, TestingValidatorCommitteeIndex)
		}
		if ret == nil {
			ret = duty
		} else {
			ret.ValidatorDuties = append(ret.ValidatorDuties, duty.ValidatorDuties...)
		}
	}

	return ret
}

func TestingCommitteeDutyWithParams(slot phase0.Slot, attestationValidatorIds []int, syncCommitteeValidatorIds []int,
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

func TestingCommitteeAttesterDuty(slot phase0.Slot, validatorIds []int) *types.CommitteeDuty {
	return TestingCommitteeDuty(slot, validatorIds, nil)
}

func TestingCommitteeSyncCommitteeDuty(slot phase0.Slot, validatorIds []int) *types.CommitteeDuty {
	return TestingCommitteeDuty(slot, nil, validatorIds)
}

// Committee duty - Attestation only
var TestingAttesterDuty = TestingCommitteeAttesterDuty(TestingDutySlot, []int{TestingValidatorIndex})

var TestingAttesterDutyNextEpoch = TestingCommitteeAttesterDuty(TestingDutySlot2, []int{TestingValidatorIndex})

var TestingAttesterDutyFirstSlot = TestingCommitteeAttesterDuty(0, []int{TestingValidatorIndex})

// Committee duty - Sync Committee only

var TestingSyncCommitteeDuty = TestingCommitteeSyncCommitteeDuty(TestingDutySlot, []int{TestingValidatorIndex})

var TestingSyncCommitteeDutyNextEpoch = TestingCommitteeSyncCommitteeDuty(TestingDutySlot2, []int{TestingValidatorIndex})

var TestingSyncCommitteeDutyFirstSlot = TestingCommitteeSyncCommitteeDuty(0, []int{TestingValidatorIndex})

// Committee duty - Attestation and Sync Committee

var TestingAttesterAndSyncCommitteeDuties = TestingCommitteeDuty(TestingDutySlot, []int{TestingValidatorIndex}, []int{TestingValidatorIndex})

var TestingAttesterAndSyncCommitteeDutiesNextEpoch = TestingCommitteeDuty(TestingDutySlot2, []int{TestingValidatorIndex}, []int{TestingValidatorIndex})

var TestingAttesterAndSyncCommitteeDutiesFirstSlot = TestingCommitteeDuty(0, []int{TestingValidatorIndex}, []int{TestingValidatorIndex})

var TestingAggregatorDutyFirstSlot = types.ValidatorDuty{
	Type:                    types.BNRoleAggregator,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    0,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingAggregatorDuty = types.ValidatorDuty{
	Type:                    types.BNRoleAggregator,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

// TestingAggregatorDutyNextEpoch testing for a second duty start
var TestingAggregatorDutyNextEpoch = types.ValidatorDuty{
	Type:                    types.BNRoleAggregator,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    TestingDutySlot2,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

// TestingSyncCommitteeContributionDutyFirstSlot
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

// TestingSyncCommitteeContributionNexEpochDuty testing for a second duty start
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

var TestingValidatorRegistrationDuty = types.ValidatorDuty{
	Type:           types.BNRoleValidatorRegistration,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot,
	ValidatorIndex: TestingValidatorIndex,
}

var TestingValidatorRegistrationDutyNextEpoch = types.ValidatorDuty{
	Type:           types.BNRoleValidatorRegistration,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot2,
	ValidatorIndex: TestingValidatorIndex,
}

var TestingVoluntaryExitDuty = types.ValidatorDuty{
	Type:           types.BNRoleVoluntaryExit,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot,
	ValidatorIndex: TestingValidatorIndex,
}

var TestingVoluntaryExitDutyNextEpoch = types.ValidatorDuty{
	Type:           types.BNRoleVoluntaryExit,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot2,
	ValidatorIndex: TestingValidatorIndex,
}

var TestingUnknownDutyType = types.ValidatorDuty{
	Type:                    UnknownDutyType,
	PubKey:                  TestingValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          22,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

var TestingWrongDutyPK = types.ValidatorDuty{
	Type:                    types.BNRoleAttester,
	PubKey:                  TestingWrongValidatorPubKey,
	Slot:                    12,
	ValidatorIndex:          TestingValidatorIndex,
	CommitteeIndex:          3,
	CommitteesAtSlot:        36,
	CommitteeLength:         128,
	ValidatorCommitteeIndex: 11,
}

//func blsSigFromHex(str string) phase0.BLSSignature {
//	byts, _ := hex.DecodeString(str)
//	ret := phase0.BLSSignature{}
//	copy(ret[:], byts)
//	return ret
//}

type TestingBeaconNode struct {
	BroadcastedRoots             []phase0.Root
	syncCommitteeAggregatorRoots map[string]bool
}

func NewTestingBeaconNode() *TestingBeaconNode {
	return &TestingBeaconNode{
		BroadcastedRoots: []phase0.Root{},
	}
}

// SetSyncCommitteeAggregatorRootHexes FOR TESTING ONLY!! sets which sync committee aggregator roots will return true for aggregator
func (bn *TestingBeaconNode) SetSyncCommitteeAggregatorRootHexes(roots map[string]bool) {
	bn.syncCommitteeAggregatorRoots = roots
}

// GetBeaconNetwork returns the beacon network the node is on
func (bn *TestingBeaconNode) GetBeaconNetwork() types.BeaconNetwork {
	return types.BeaconTestNetwork
}

// GetAttestationData returns attestation data by the given slot and committee index
func (bn *TestingBeaconNode) GetAttestationData(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) (*phase0.
	AttestationData, spec.DataVersion, error) {
	data := *TestingAttestationData
	data.Slot = slot
	if committeeIndex != 0 {
		data.Index = committeeIndex
	}
	return &data, spec.DataVersionPhase0, nil
}

// SubmitAttestations submit attestations to the node
func (bn *TestingBeaconNode) SubmitAttestations(attestations []*phase0.Attestation) error {
	for _, att := range attestations {
		r, _ := att.HashTreeRoot()
		bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	}
	return nil
}

func (bn *TestingBeaconNode) SubmitValidatorRegistration(pubkey []byte, feeRecipient bellatrix.ExecutionAddress, sig phase0.BLSSignature) error {
	pk := phase0.BLSPubKey{}
	copy(pk[:], pubkey)

	vr := v1.ValidatorRegistration{
		FeeRecipient: feeRecipient,
		GasLimit:     TestingValidatorRegistration.GasLimit,
		Timestamp:    TestingValidatorRegistration.Timestamp,
		Pubkey:       pk,
	}

	r, _ := vr.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// SubmitVoluntaryExit submit the VoluntaryExit object to the node
func (bn *TestingBeaconNode) SubmitVoluntaryExit(voluntaryExit *phase0.SignedVoluntaryExit) error {
	r, _ := voluntaryExit.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// GetBeaconBlock returns beacon block by the given slot, graffiti, and randao.
func (bn *TestingBeaconNode) GetBeaconBlock(slot phase0.Slot, graffiti, randao []byte) (ssz.Marshaler, spec.DataVersion, error) {
	version := VersionBySlot(slot)
	vBlk := TestingBeaconBlockV(version)

	switch version {
	case spec.DataVersionCapella:
		return vBlk.Capella, version, nil
	case spec.DataVersionDeneb:
		return vBlk.Deneb, version, nil
	default:
		panic("unsupported version")
	}
}

// SubmitBeaconBlock submit the block to the node
func (bn *TestingBeaconNode) SubmitBeaconBlock(block *api.VersionedProposal, sig phase0.BLSSignature) error {
	var r [32]byte

	switch block.Version {
	case spec.DataVersionCapella:
		if block.Capella == nil {
			return errors.Errorf("%s block is nil", block.Version.String())
		}
		sb := &capella.SignedBeaconBlock{
			Message:   block.Capella,
			Signature: sig,
		}
		r, _ = sb.HashTreeRoot()
	case spec.DataVersionDeneb:
		if block.Deneb == nil {
			return errors.Errorf("%s block contents is nil", block.Version.String())
		}
		if block.Deneb.Block == nil {
			return errors.Errorf("%s block is nil", block.Version.String())
		}
		sb := &apiv1deneb.SignedBlockContents{
			SignedBlock: &deneb.SignedBeaconBlock{
				Message:   block.Deneb.Block,
				Signature: sig,
			},
			KZGProofs: block.Deneb.KZGProofs,
			Blobs:     block.Deneb.Blobs,
		}
		r, _ = sb.HashTreeRoot()
	default:
		return errors.Errorf("unknown block version %d", block.Version)
	}

	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// SubmitBlindedBeaconBlock submit the blinded block to the node
func (bn *TestingBeaconNode) SubmitBlindedBeaconBlock(block *api.VersionedBlindedProposal, sig phase0.BLSSignature) error {
	var r [32]byte

	switch block.Version {
	case spec.DataVersionCapella:
		if block.Capella == nil {
			return errors.Errorf("%s blinded block is nil", block.Version.String())
		}
		sb := &apiv1capella.SignedBlindedBeaconBlock{
			Message:   block.Capella,
			Signature: sig,
		}
		r, _ = sb.HashTreeRoot()
	case spec.DataVersionDeneb:
		if block.Deneb == nil {
			return errors.Errorf("%s blinded block is nil", block.Version.String())
		}
		sb := &apiv1deneb.SignedBlindedBeaconBlock{
			Message:   block.Deneb,
			Signature: sig,
		}
		r, _ = sb.HashTreeRoot()
	default:
		return errors.Errorf("unknown blinded block version %s", block.Version.String())
	}

	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// SubmitAggregateSelectionProof returns an AggregateAndProof object
func (bn *TestingBeaconNode) SubmitAggregateSelectionProof(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, committeeLength uint64, index phase0.ValidatorIndex, slotSig []byte) (ssz.Marshaler, spec.DataVersion, error) {
	return TestingAggregateAndProof, spec.DataVersionPhase0, nil
}

// SubmitSignedAggregateSelectionProof broadcasts a signed aggregator msg
func (bn *TestingBeaconNode) SubmitSignedAggregateSelectionProof(msg *phase0.SignedAggregateAndProof) error {
	r, _ := msg.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// GetSyncMessageBlockRoot returns beacon block root for sync committee
func (bn *TestingBeaconNode) GetSyncMessageBlockRoot(slot phase0.Slot) (phase0.Root, spec.DataVersion, error) {
	return TestingSyncCommitteeBlockRoot, spec.DataVersionPhase0, nil
}

// SubmitSyncMessage submits a signed sync committee msg
func (bn *TestingBeaconNode) SubmitSyncMessages(msgs []*altair.SyncCommitteeMessage) error {
	for _, msg := range msgs {
		r, _ := msg.HashTreeRoot()
		bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	}
	return nil
}

// IsSyncCommitteeAggregator returns tru if aggregator
func (bn *TestingBeaconNode) IsSyncCommitteeAggregator(proof []byte) (bool, error) {
	if len(bn.syncCommitteeAggregatorRoots) != 0 {
		if val, found := bn.syncCommitteeAggregatorRoots[hex.EncodeToString(proof)]; found {
			return val, nil
		}
		return false, nil
	}
	return true, nil
}

// SyncCommitteeSubnetID returns sync committee subnet ID from subcommittee index
func (bn *TestingBeaconNode) SyncCommitteeSubnetID(index phase0.CommitteeIndex) (uint64, error) {
	// each subcommittee index correlates to TestingContributionProofRoots by index
	return uint64(index), nil
}

// GetSyncCommitteeContribution returns
func (bn *TestingBeaconNode) GetSyncCommitteeContribution(slot phase0.Slot, selectionProofs []phase0.BLSSignature, subnetIDs []uint64) (ssz.Marshaler, spec.DataVersion, error) {
	return &TestingContributionsData, spec.DataVersionBellatrix, nil
}

// SubmitSignedContributionAndProof broadcasts to the network
func (bn *TestingBeaconNode) SubmitSignedContributionAndProof(contribution *altair.SignedContributionAndProof) error {
	r, _ := contribution.HashTreeRoot()
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

func (bn *TestingBeaconNode) DomainData(epoch phase0.Epoch, domain phase0.DomainType) (phase0.Domain, error) {
	// epoch is used to calculate fork version, here we hard code it
	return types.ComputeETHDomain(domain, types.GenesisForkVersion, types.GenesisValidatorsRoot)
}
