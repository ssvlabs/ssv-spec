package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/types"
)

func TestAggregatorCommitteeConsensusData(version spec.DataVersion) *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(TestingAggregatorCommitteeDutyMixed(version), version, nil)
}

// TestAggregatorCommitteeConsensusDataForDuty creates consensus data matching the given duty
func TestAggregatorCommitteeConsensusDataForDuty(duty *types.AggregatorCommitteeDuty, version spec.DataVersion, ksMap map[phase0.ValidatorIndex]*TestKeySet) *types.AggregatorCommitteeConsensusData {
	consensusData := &types.AggregatorCommitteeConsensusData{
		Version: version,

		// Initialize empty slices
		Aggregators:                 []types.AssignedAggregator{},
		AggregatorsCommitteeIndexes: []uint64{},
		AggregatedAttestations:      [][]byte{},
		Contributors:                []types.AssignedAggregator{},
		SyncCommitteeContributions:  []altair.SyncCommitteeContribution{},
	}

	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	aggDomainData, _ := beacon.DomainData(1, types.DomainSelectionProof)
	sccDomainData, _ := beacon.DomainData(1, types.DomainSyncCommitteeSelectionProof)
	signingEnabled := (len(ksMap) > 0)

	// Process each validator duty
	for _, validatorDuty := range duty.ValidatorDuties {
		if validatorDuty == nil {
			continue
		}

		switch validatorDuty.Type {
		case types.BNRoleAggregator:
			// Create attestation for this validator based on version
			var marshaledAtt []byte

			if version >= spec.DataVersionElectra {
				// For Electra and newer versions, create an electra.Attestation with CommitteeBits
				attestation := &electra.Attestation{
					AggregationBits: bitfield.NewBitlist(128),
					Signature:       phase0.BLSSignature{},
					Data:            TestingAttestationData(version),
					CommitteeBits:   bitfield.NewBitvector64(),
				}
				// Leave AggregationBits empty for testing
				// Leave CommitteeBits empty for testing - in reality they would be set
				marshaledAtt, _ = attestation.MarshalSSZ()
			} else {
				// For pre-Electra, create a phase0.Attestation
				attestation := &phase0.Attestation{
					AggregationBits: bitfield.NewBitlist(128),
					Signature:       phase0.BLSSignature{},
					Data:            TestingAttestationData(version),
				}
				// Leave AggregationBits empty for testing
				marshaledAtt, _ = attestation.MarshalSSZ()
			}

			blsSig := phase0.BLSSignature{}
			if signingEnabled {
				blsSignature, _, _ := signer.SignBeaconObject(types.SSZUint64(duty.DutySlot()), aggDomainData, ksMap[validatorDuty.ValidatorIndex].ValidatorSK.GetPublicKey().Serialize(), types.DomainSelectionProof)
				blsSig = phase0.BLSSignature(blsSignature)
			}

			// Add aggregator data
			consensusData.Aggregators = append(consensusData.Aggregators, types.AssignedAggregator{
				ValidatorIndex: validatorDuty.ValidatorIndex,
				CommitteeIndex: uint64(validatorDuty.CommitteeIndex),
				SelectionProof: blsSig,
			})

			commIndexAlreadyExists := false
			for _, commIndex := range consensusData.AggregatorsCommitteeIndexes {
				if commIndex == uint64(validatorDuty.CommitteeIndex) {
					commIndexAlreadyExists = true
					break
				}
			}
			if !commIndexAlreadyExists {
				consensusData.AggregatorsCommitteeIndexes = append(consensusData.AggregatorsCommitteeIndexes, uint64(validatorDuty.CommitteeIndex))
				consensusData.AggregatedAttestations = append(consensusData.AggregatedAttestations, marshaledAtt)
			}

		case types.BNRoleSyncCommitteeContribution:
			for idx, validatorSyncCommitteeIndex := range validatorDuty.ValidatorSyncCommitteeIndices {

				subnet := beacon.SyncCommitteeSubnetID(phase0.CommitteeIndex(validatorSyncCommitteeIndex))

				data := &altair.SyncAggregatorSelectionData{
					Slot:              duty.DutySlot(),
					SubcommitteeIndex: subnet,
				}

				blsSig := TestingContributionProofsSigned[idx]
				if signingEnabled {
					blsSignature, _, _ := signer.SignBeaconObject(data, sccDomainData, ksMap[validatorDuty.ValidatorIndex].ValidatorSK.GetPublicKey().Serialize(), types.DomainSyncCommitteeSelectionProof)
					blsSig = phase0.BLSSignature(blsSignature)
				}

				contribution := TestingSyncCommitteeContributions[subnet]

				// Add sync committee contributor data
				consensusData.Contributors = append(consensusData.Contributors, types.AssignedAggregator{
					ValidatorIndex: validatorDuty.ValidatorIndex,
					SelectionProof: blsSig,
					CommitteeIndex: contribution.SubcommitteeIndex,
				})

				commIndexAlreadyExists := false
				for _, commIndex := range consensusData.SyncCommitteeContributions {
					if commIndex.SubcommitteeIndex == contribution.SubcommitteeIndex {
						commIndexAlreadyExists = true
						break
					}
				}
				if !commIndexAlreadyExists {
					consensusData.SyncCommitteeContributions = append(consensusData.SyncCommitteeContributions, *contribution)
				}
			}
		}
	}

	return consensusData
}

// TestAggregatorCommitteeConsensusDataBytesForDuty encodes the consensus data for the given duty
func TestAggregatorCommitteeConsensusDataBytesForDuty(duty *types.AggregatorCommitteeDuty, version spec.DataVersion) []byte {
	consensusData := TestAggregatorCommitteeConsensusDataForDuty(duty, version, nil)
	bytes, _ := consensusData.Encode()
	return bytes
}

// TestAggregatorCommitteeConsensusDataBytesForDuty encodes the consensus data for the given duty
func TestAggregatorCommitteeConsensusDataBytesForDutyWithKS(duty *types.AggregatorCommitteeDuty, version spec.DataVersion, ksMap map[phase0.ValidatorIndex]*TestKeySet) []byte {
	consensusData := TestAggregatorCommitteeConsensusDataForDuty(duty, version, ksMap)
	bytes, _ := consensusData.Encode()
	return bytes
}
