package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Versioned AggregatorCommitteeDuty
// ==================================================

func TestingAggregatorCommitteeDutyForSlot(slot phase0.Slot, aggregatorValidatorIds []int, syncCommitteeValidatorIds []int) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDutyWithParams(
		slot,
		aggregatorValidatorIds,
		syncCommitteeValidatorIds,
		TestingCommitteeIndex,
		TestingCommitteesAtSlot,
		TestingCommitteeLenght,
		TestingValidatorCommitteeIndex)
}

func TestingAggregatorCommitteeDuty(aggregatorValidatorIds []int, syncCommitteeValidatorIds []int, version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDutyForSlot(TestingDutySlotV(version), aggregatorValidatorIds, syncCommitteeValidatorIds)
}

func TestingAggregatorCommitteeDutyWithParams(
	slot phase0.Slot,
	aggregatorValidatorIds []int,
	syncCommitteeValidatorIds []int,
	committeeIndex phase0.CommitteeIndex,
	committeesAtSlot uint64,
	committeeLenght uint64,
	validatorCommitteeIndex uint64) *types.AggregatorCommitteeDuty {

	duties := make([]*types.ValidatorDuty, 0)

	for _, valIdx := range aggregatorValidatorIds {
		pk := getValPubKeyByValIdx(valIdx)
		duties = append(duties, &types.ValidatorDuty{
			Type:                    types.BNRoleAggregator,
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
			Type:                          types.BNRoleSyncCommitteeContribution,
			PubKey:                        pk,
			Slot:                          slot,
			ValidatorIndex:                phase0.ValidatorIndex(valIdx),
			CommitteeIndex:                3,
			CommitteesAtSlot:              committeesAtSlot,
			CommitteeLength:               committeeLenght,
			ValidatorCommitteeIndex:       validatorCommitteeIndex,
			ValidatorSyncCommitteeIndices: TestingContributionProofIndexes,
		})
	}

	return &types.AggregatorCommitteeDuty{
		Slot:            slot,
		ValidatorDuties: duties,
	}
}

// ==================================================
// AggregatorCommitteeDuty - Aggregator only
// ==================================================

var TestingAggregatorCommitteeDutySingle = func(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDuty([]int{TestingValidatorIndex}, nil, version)
}

var TestingAggregatorDutyForValidator = func(version spec.DataVersion, validatorIndex phase0.ValidatorIndex) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDuty([]int{int(validatorIndex)}, nil, version)
}

var TestingAggregatorDutyForValidators = func(version spec.DataVersion, validatorIndexLst []int) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDuty(validatorIndexLst, nil, version)
}

// ==================================================
// AggregatorCommitteeDuty - SyncCommittee Contributor only
// ==================================================

var TestingSyncCommitteeContributorDuty = func(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDuty(nil, []int{TestingValidatorIndex}, version)
}

var TestingSyncCommitteeContributorDutyForValidators = func(version spec.DataVersion, validatorIndexList []int) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDuty(nil, validatorIndexList, version)
}

// ==================================================
// AggregatorCommitteeDuty - Aggregator and SyncCommittee Contributor
// ==================================================

var TestingAggregatorAndSyncCommitteeContributorDuties = func(version spec.DataVersion) *types.AggregatorCommitteeDuty {
	return TestingAggregatorCommitteeDuty([]int{TestingValidatorIndex}, []int{TestingValidatorIndex}, version)
}

// ==================================================
// Beacon Roots for AggregatorCommittee duty
// ==================================================

var TestingSignedAggregatorCommitteeBeaconObjectSSZRoot = func(duty *types.AggregatorCommitteeDuty, ksMap map[phase0.ValidatorIndex]*TestKeySet, version spec.DataVersion) []string {
	ret := make([]string, 0)

	// Create consensus data that matches what the runner will use
	consensusData := TestAggregatorCommitteeConsensusDataForDuty(duty, version)

	// Get aggregate and proofs from consensus data
	aggregateAndProofs, _, err := consensusData.GetAggregateAndProofs()
	if err != nil {
		panic(err)
	}

	for i, aggregateAndProof := range aggregateAndProofs {
		validatorIndex := consensusData.Aggregators[i].ValidatorIndex
		ks := ksMap[validatorIndex]
		if ks == nil {
			continue
		}

		// Sign the aggregate and proof
		signer := NewTestingKeyManager()
		beacon := NewTestingBeaconNode()
		d, _ := beacon.DomainData(1, types.DomainAggregateAndProof)

		// Get the appropriate aggregate and proof object
		var signingRoot ssz.HashRoot
		switch version {
		case spec.DataVersionElectra:
			signingRoot = aggregateAndProof.Electra
		default:
			// Get the appropriate version field
			switch aggregateAndProof.Version {
			case spec.DataVersionPhase0:
				signingRoot = aggregateAndProof.Phase0
			case spec.DataVersionAltair:
				signingRoot = aggregateAndProof.Altair
			case spec.DataVersionBellatrix:
				signingRoot = aggregateAndProof.Bellatrix
			case spec.DataVersionCapella:
				signingRoot = aggregateAndProof.Capella
			case spec.DataVersionDeneb:
				signingRoot = aggregateAndProof.Deneb
			}
		}

		sig, _, _ := signer.SignBeaconObject(signingRoot, d, ks.ValidatorPK.Serialize(), types.DomainAggregateAndProof)

		// Convert signature to BLSSignature
		var blsSig phase0.BLSSignature
		copy(blsSig[:], sig)

		// Create signed aggregate and proof
		var signedAgg ssz.HashRoot
		switch version {
		case spec.DataVersionElectra:
			signedAgg = &electra.SignedAggregateAndProof{
				Message:   aggregateAndProof.Electra,
				Signature: blsSig,
			}
		default:
			// For pre-electra versions, use phase0
			signedAgg = &phase0.SignedAggregateAndProof{
				Message:   aggregateAndProof.Phase0,
				Signature: blsSig,
			}
		}

		ret = append(ret, GetSSZRootNoError(signedAgg))
	}

	// TODO: Handle sync committee contributions

	return ret
}
