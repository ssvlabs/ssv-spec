package maxmsgsize

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	maxSizeAggregatorCommitteeConsensusData = 8970524
	maxSizePhase0Attestation                = 2276
	maxSizeElectraAttestation               = 131308
)

func maxAggregatorCommitteeConsensusData() *types.AggregatorCommitteeConsensusData {

	maxVals := 3000

	// Aggregator fields
	aggs := make([]types.AssignedAggregator, 0)
	for i := 0; i < maxVals; i++ {
		aggs = append(aggs, types.AssignedAggregator{
			ValidatorIndex: phase0.ValidatorIndex(i),
			SelectionProof: phase0.BLSSignature([96]byte{1}),
			CommitteeIndex: uint64(i),
		})
	}
	maxCommIdxs := 64
	aggCommIdxs := make([]uint64, 0)
	for i := 0; i < maxCommIdxs; i++ {
		aggCommIdxs = append(aggCommIdxs, 1)
	}
	aggAtt := make([][]byte, 0)
	for i := 0; i < maxCommIdxs; i++ {
		att := [131308]byte{1}
		aggAtt = append(aggAtt, att[:])
	}

	// Contributor fields
	maxContribs := 2048
	contributors := make([]types.AssignedAggregator, 0)
	for i := 0; i < maxContribs; i++ {
		contributors = append(contributors, types.AssignedAggregator{
			ValidatorIndex: phase0.ValidatorIndex(i),
			SelectionProof: phase0.BLSSignature([96]byte{1}),
			CommitteeIndex: uint64(i),
		})
	}

	maxSCC := 4
	syncCommitteeContributions := make([]altair.SyncCommitteeContribution, 0)
	bitvec := bitfield.NewBitvector128()
	for i := 0; i < maxSCC; i++ {
		syncCommitteeContributions = append(syncCommitteeContributions, altair.SyncCommitteeContribution{
			Slot:              1,
			BeaconBlockRoot:   [32]byte{1},
			SubcommitteeIndex: 1,
			AggregationBits:   bitvec,
			Signature:         phase0.BLSSignature([96]byte{1}),
		})
	}

	return &types.AggregatorCommitteeConsensusData{
		Version: spec.DataVersionPhase0,

		Aggregators:                 aggs,
		AggregatorsCommitteeIndexes: aggCommIdxs,
		AggregatedAttestations:      aggAtt,

		Contributors:               contributors,
		SyncCommitteeContributions: syncCommitteeContributions,
	}
}

// Wrapper types to allow defining Encode/Decode methods on external types
type Phase0AttestationWrapper struct {
	*phase0.Attestation
}

func (w *Phase0AttestationWrapper) Encode() ([]byte, error) {
	return w.Attestation.MarshalSSZ()
}

func (w *Phase0AttestationWrapper) Decode(data []byte) error {
	return w.Attestation.UnmarshalSSZ(data)
}

type ElectraAttestationWrapper struct {
	*electra.Attestation
}

func (w *ElectraAttestationWrapper) Encode() ([]byte, error) {
	return w.Attestation.MarshalSSZ()
}

func (w *ElectraAttestationWrapper) Decode(data []byte) error {
	return w.Attestation.UnmarshalSSZ(data)
}

func maxPhase0Attestation() *Phase0AttestationWrapper {
	aggbits := [2048]byte{1}
	return &Phase0AttestationWrapper{
		Attestation: &phase0.Attestation{
			AggregationBits: bitfield.Bitlist(aggbits[:]),
			Data: &phase0.AttestationData{
				Slot:            1,
				Index:           0,
				BeaconBlockRoot: [32]byte{1},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
				Target: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
			},
			Signature: phase0.BLSSignature([96]byte{1}),
		},
	}
}

func maxElectraAttestation() *ElectraAttestationWrapper {
	aggbits := [131072]byte{1}
	return &ElectraAttestationWrapper{
		Attestation: &electra.Attestation{
			AggregationBits: bitfield.Bitlist(aggbits[:]),
			Data: &phase0.AttestationData{
				Slot:            1,
				Index:           0,
				BeaconBlockRoot: [32]byte{1},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
				Target: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
			},
			Signature:     phase0.BLSSignature([96]byte{1}),
			CommitteeBits: bitfield.NewBitvector64(),
		},
	}
}

func MaxAggregatorCommitteeConsensusData() *StructureSizeTest {
	return NewStructureSizeTest(
		"max AggregatorCommitteeConsensusData",
		testdoc.StructureSizeTestMaxAggregatorCommitteeConsensusDataDoc,
		maxAggregatorCommitteeConsensusData(),
		maxSizeAggregatorCommitteeConsensusData,
		true,
	)
}

func MaxPhase0Attestation() *StructureSizeTest {
	return NewStructureSizeTest(
		"max Phase0Attestation",
		testdoc.StructureSizeTestMaxPhase0AttestationDoc,
		maxPhase0Attestation(),
		maxSizePhase0Attestation,
		true,
	)
}

func MaxElectraAttestation() *StructureSizeTest {
	return NewStructureSizeTest(
		"max ElectraAttestation",
		testdoc.StructureSizeTestMaxElectraAttestationDoc,
		maxElectraAttestation(),
		maxSizeElectraAttestation,
		true,
	)
}
