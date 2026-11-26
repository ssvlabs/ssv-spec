package maxmsgsize

import (
	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/altair"
	eth2gloas "github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	MaxSizeAggregatorCommitteeConsensusData = 8970524
	// MaxSizeAggregatedAttestation is the bound on one AggregatedAttestations element (consensus_data.go's
	// ssz-max:"64,131308"), and the largest encoded aggregate the decoder accepts. It is a root-affecting limit,
	// kept well above the largest attestation the beacon chain produces, which is 16621 bytes
	// (MaxGloasAttestation). Do not "correct" it down.
	MaxSizeAggregatedAttestation = 131308
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
		att := [MaxSizeAggregatedAttestation]byte{1}
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
		Version: gloas.DataVersionGloas,

		Aggregators:                 aggs,
		AggregatorsCommitteeIndexes: aggCommIdxs,
		AggregatedAttestations:      aggAtt,

		Contributors:               contributors,
		SyncCommitteeContributions: syncCommitteeContributions,
	}
}

// GloasAttestationWrapper wraps go-eth2-client's Gloas attestation, the container the spec decodes an aggregate
// as (SIP #94 §2), to give it Encode/Decode methods.
type GloasAttestationWrapper struct {
	*eth2gloas.Attestation
}

func (w *GloasAttestationWrapper) Encode() ([]byte, error) {
	return w.MarshalSSZ()
}

func (w *GloasAttestationWrapper) Decode(data []byte) error {
	return w.UnmarshalSSZ(data)
}

// maxBeaconChainAggregationBits is the most aggregation bits an attestation has on the beacon chain:
// MAX_VALIDATORS_PER_COMMITTEE * MAX_COMMITTEES_PER_SLOT.
const maxBeaconChainAggregationBits = 131072

func maxGloasAttestation() *GloasAttestationWrapper {
	return &GloasAttestationWrapper{
		Attestation: &eth2gloas.Attestation{
			AggregationBits: bitfield.NewBitlist(maxBeaconChainAggregationBits),
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
		MaxSizeAggregatorCommitteeConsensusData,
		true,
	)
}

// MaxGloasAttestation pins the encoded size of the largest attestation the beacon chain produces: one with
// aggregation bits for every validator of every committee of a slot. It is the beacon-chain maximum, not a bound
// of the container: the Gloas aggregation bits are a progressive bitlist, and the spec accepts an aggregate up
// to MaxSizeAggregatedAttestation bytes.
func MaxGloasAttestation() *StructureSizeTest {
	return NewStructureSizeTest(
		"max GloasAttestation",
		testdoc.StructureSizeTestMaxGloasAttestationDoc,
		maxGloasAttestation(),
		16621,
		// The aggregation bits have no ssz-max tag to be equal to, so the tag check has nothing to hold them to.
		false,
	)
}
