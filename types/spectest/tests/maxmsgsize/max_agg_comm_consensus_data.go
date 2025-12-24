package maxmsgsize

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	maxSizeAggregatorCommitteeConsensusData = 67675676
)

func maxAggregatorCommitteeConsenssuData() *types.AggregatorCommitteeConsensusData {

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
		att := [1048576]byte{1}
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

func MaxAggregatorCommitteeConsensusData() *StructureSizeTest {
	return NewStructureSizeTest(
		"max AggregatorCommitteeConsensusData",
		testdoc.StructureSizeTestMaxAggregatorCommitteeConsensusDataDoc,
		maxAggregatorCommitteeConsenssuData(),
		maxSizeAggregatorCommitteeConsensusData,
		true,
	)
}
