package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Aggregator
// ==================================================

var TestAggregatorConsensusData = func(version spec.DataVersion) *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(TestingAggregatorCommitteeDutyOnlyAggregator(version), version, nil)
}
var TestAggregatorConsensusDataByts = func(version spec.DataVersion) []byte {
	byts, _ := TestAggregatorConsensusData(version).Encode()
	return byts
}

// ==================================================
// Attester
// ==================================================

// Used only as invalid test case
var TestAttesterConsensusData = &types.ProposerConsensusData{
	Duty:    *TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0],
	DataSSZ: TestingAttestationDataBytes(spec.DataVersionPhase0),
	Version: spec.DataVersionPhase0,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

// ==================================================
// Sync Committee
// ==================================================

// Used only as invalid test case
var TestSyncCommitteeConsensusData = &types.ProposerConsensusData{
	Duty:    *TestingSyncCommitteeDuty(spec.DataVersionPhase0).ValidatorDuties[0],
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
	Version: spec.DataVersionPhase0,
}
var TestSyncCommitteeConsensusDataByts, _ = TestSyncCommitteeConsensusData.Encode()

// ==================================================
// Proposer
// ==================================================

var TestProposerConsensusDataV = func(version spec.DataVersion) *types.ProposerConsensusData {
	duty := TestingProposerDutyV(version)
	return &types.ProposerConsensusData{
		Duty:    *duty,
		Version: version,
		DataSSZ: TestingBeaconBlockBytesV(version),
	}
}

var TestProposerConsensusDataBytsV = func(version spec.DataVersion) []byte {
	cd := TestProposerConsensusDataV(version)
	byts, _ := cd.Encode()
	return byts
}

var TestProposerBlindedBlockConsensusDataV = func(version spec.DataVersion) *types.ProposerConsensusData {
	return &types.ProposerConsensusData{
		Duty:    *TestingProposerDutyV(version),
		Version: version,
		DataSSZ: TestingBlindedBeaconBlockBytesV(version),
	}
}

var TestProposerBlindedBlockConsensusDataBytsV = func(version spec.DataVersion) []byte {
	cd := TestProposerBlindedBlockConsensusDataV(version)
	byts, _ := cd.Encode()
	return byts
}

// TestProposerConsensusDataMutatedBlockByts returns cd encoded with its block's slot and proposer index changed by
// mutate. cd's own duty is encoded as is.
func TestProposerConsensusDataMutatedBlockByts(cd *types.ProposerConsensusData, mutate func(*phase0.Slot, *phase0.ValidatorIndex)) []byte {
	blk, _, err := cd.GetBlockData() // decodes a fresh copy
	if err != nil {
		panic(err.Error())
	}

	var obj ssz.Marshaler
	switch {
	case blk.Version == spec.DataVersionCapella && blk.Blinded:
		mutate(&blk.CapellaBlinded.Slot, &blk.CapellaBlinded.ProposerIndex)
		obj = blk.CapellaBlinded
	case blk.Version == spec.DataVersionCapella:
		mutate(&blk.Capella.Slot, &blk.Capella.ProposerIndex)
		obj = blk.Capella
	case blk.Version == spec.DataVersionDeneb && blk.Blinded:
		mutate(&blk.DenebBlinded.Slot, &blk.DenebBlinded.ProposerIndex)
		obj = blk.DenebBlinded
	case blk.Version == spec.DataVersionDeneb:
		mutate(&blk.Deneb.Block.Slot, &blk.Deneb.Block.ProposerIndex)
		obj = blk.Deneb
	case blk.Version == spec.DataVersionElectra && blk.Blinded:
		mutate(&blk.ElectraBlinded.Slot, &blk.ElectraBlinded.ProposerIndex)
		obj = blk.ElectraBlinded
	case blk.Version == spec.DataVersionElectra:
		mutate(&blk.Electra.Block.Slot, &blk.Electra.Block.ProposerIndex)
		obj = blk.Electra
	case blk.Version == spec.DataVersionFulu && blk.Blinded:
		mutate(&blk.FuluBlinded.Slot, &blk.FuluBlinded.ProposerIndex)
		obj = blk.FuluBlinded
	case blk.Version == spec.DataVersionFulu:
		mutate(&blk.Fulu.Block.Slot, &blk.Fulu.Block.ProposerIndex)
		obj = blk.Fulu
	default:
		panic("unsupported version")
	}

	dataSSZ, err := obj.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	ret, err := (&types.ProposerConsensusData{Duty: cd.Duty, Version: cd.Version, DataSSZ: dataSSZ}).Encode()
	if err != nil {
		panic(err.Error())
	}
	return ret
}

// ==================================================
// Sync Committee Contribution
// ==================================================

var TestSyncCommitteeContributionConsensusDataF = func() *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(TestingAggregatorCommitteeDutyOnlySyncCommittee(), spec.DataVersionPhase0, nil)
}

var TestSyncCommitteeContributionConsensusDataForDuty = func(duty *types.AggregatorCommitteeDuty) *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(duty, spec.DataVersionPhase0, nil)
}

var TestSyncCommitteeContributionConsensusData = TestSyncCommitteeContributionConsensusDataF()

var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()
