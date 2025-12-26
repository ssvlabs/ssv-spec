package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Aggregator
// ==================================================

var TestAggregatorConsensusData = func(version spec.DataVersion) *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(TestingAggregatorCommitteeDutyOnlyAggregator(version), version)
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

// ==================================================
// Sync Committee Contribution
// ==================================================

var TestSyncCommitteeContributionConsensusDataF = func() *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(TestingAggregatorCommitteeDutyOnlySyncCommittee(), spec.DataVersionPhase0)
}

var TestSyncCommitteeContributionConsensusDataForDuty = func(duty *types.AggregatorCommitteeDuty) *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(duty, spec.DataVersionPhase0)
}

var TestSyncCommitteeContributionConsensusData = TestSyncCommitteeContributionConsensusDataF()

var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()
