package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
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
	Duty:    *TestingAttesterDuty(gloas.DataVersionGloas).ValidatorDuties[0],
	DataSSZ: TestingAttestationDataBytes(gloas.DataVersionGloas),
	Version: gloas.DataVersionGloas,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

// ==================================================
// Sync Committee
// ==================================================

// Used only as invalid test case
var TestSyncCommitteeConsensusData = &types.ProposerConsensusData{
	Duty:    *TestingSyncCommitteeDuty(gloas.DataVersionGloas).ValidatorDuties[0],
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
	Version: gloas.DataVersionGloas,
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

// ==================================================
// Sync Committee Contribution
// ==================================================

var TestSyncCommitteeContributionConsensusDataF = func() *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(TestingAggregatorCommitteeDutyOnlySyncCommittee(), gloas.DataVersionGloas, nil)
}

var TestSyncCommitteeContributionConsensusDataForDuty = func(duty *types.AggregatorCommitteeDuty) *types.AggregatorCommitteeConsensusData {
	return TestAggregatorCommitteeConsensusDataForDuty(duty, gloas.DataVersionGloas, nil)
}

var TestSyncCommitteeContributionConsensusData = TestSyncCommitteeContributionConsensusDataF()

var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()
