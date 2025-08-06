package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/types"
)

// TestAggregatorCommitteeConsensusDataForDuty creates consensus data matching the given duty
func TestAggregatorCommitteeConsensusDataForDuty(duty *types.AggregatorCommitteeDuty, version spec.DataVersion) *types.AggregatorCommitteeConsensusData {
	consensusData := &types.AggregatorCommitteeConsensusData{
		Version: version,

		// Initialize empty slices
		Aggregators:                 []types.AssignedAggregator{},
		AggregatorsCommitteeIndexes: []uint64{},
		Attestations:                [][]byte{},
		Contributors:                []types.AssignedAggregator{},
		SyncCommitteeSubnets:        []uint64{},
		SyncCommitteeContributions:  []altair.SyncCommitteeContribution{},
	}

	// Process each validator duty
	for _, validatorDuty := range duty.ValidatorDuties {
		if validatorDuty == nil {
			continue
		}

		switch validatorDuty.Type {
		case types.BNRoleAggregator:
			// Create attestation for this validator based on version
			var marshaledAtt []byte

			if version == spec.DataVersionElectra {
				// For Electra, create an electra.Attestation with CommitteeBits
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

			// Add aggregator data
			consensusData.Aggregators = append(consensusData.Aggregators, types.AssignedAggregator{
				ValidatorIndex: validatorDuty.ValidatorIndex,
				CommitteeIndex: uint64(validatorDuty.CommitteeIndex),
			})
			consensusData.AggregatorsCommitteeIndexes = append(consensusData.AggregatorsCommitteeIndexes, uint64(validatorDuty.CommitteeIndex))
			consensusData.Attestations = append(consensusData.Attestations, marshaledAtt)

		case types.BNRoleSyncCommitteeContribution:
			for i, contribution := range TestingSyncCommitteeContributions {
				// Add sync committee contributor data
				consensusData.Contributors = append(consensusData.Contributors, types.AssignedAggregator{
					ValidatorIndex: validatorDuty.ValidatorIndex,
					SelectionProof: TestingContributionProofsSigned[i],
				})

				// Add sync committee contribution
				consensusData.SyncCommitteeContributions = append(consensusData.SyncCommitteeContributions, *contribution)

				// Append correct SubcommitteeIndex
				consensusData.SyncCommitteeSubnets = append(consensusData.SyncCommitteeSubnets, uint64(contribution.SubcommitteeIndex))
			}
		}
	}

	return consensusData
}

// TestAggregatorCommitteeConsensusDataBytesForDuty encodes the consensus data for the given duty
func TestAggregatorCommitteeConsensusDataBytesForDuty(duty *types.AggregatorCommitteeDuty, version spec.DataVersion) []byte {
	consensusData := TestAggregatorCommitteeConsensusDataForDuty(duty, version)
	bytes, _ := consensusData.Encode()
	return bytes
}
