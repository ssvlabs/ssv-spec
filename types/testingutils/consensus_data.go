package testingutils

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Aggregator
// ==================================================

var TestAggregatorConsensusData = func(version spec.DataVersion) *types.ValidatorConsensusData {
	return &types.ValidatorConsensusData{
		Duty:    *TestingAggregatorDuty(version),
		DataSSZ: TestingAggregateAndProofBytesV(version),
		Version: version,
	}
}
var TestAggregatorConsensusDataByts = func(version spec.DataVersion) []byte {
	byts, _ := TestAggregatorConsensusData(version).Encode()
	return byts
}

// ==================================================
// Attester
// ==================================================

// Used only as invalid test case
var TestAttesterConsensusData = &types.ValidatorConsensusData{
	Duty:    *TestingAttesterDuty(spec.DataVersionPhase0).ValidatorDuties[0],
	DataSSZ: TestingAttestationDataBytes(spec.DataVersionPhase0),
	Version: spec.DataVersionPhase0,
}
var TestAttesterConsensusDataByts, _ = TestAttesterConsensusData.Encode()

// ==================================================
// Sync Committee
// ==================================================

// Used only as invalid test case
var TestSyncCommitteeConsensusData = &types.ValidatorConsensusData{
	Duty:    *TestingSyncCommitteeDuty(spec.DataVersionPhase0).ValidatorDuties[0],
	DataSSZ: TestingSyncCommitteeBlockRoot[:],
	Version: spec.DataVersionPhase0,
}
var TestSyncCommitteeConsensusDataByts, _ = TestSyncCommitteeConsensusData.Encode()

// ==================================================
// Proposer
// ==================================================

var TestProposerConsensusDataV = func(version spec.DataVersion) *types.ValidatorConsensusData {
	duty := TestingProposerDutyV(version)
	return &types.ValidatorConsensusData{
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

var TestProposerBlindedBlockConsensusDataV = func(version spec.DataVersion) *types.ValidatorConsensusData {
	return &types.ValidatorConsensusData{
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

var TestSyncCommitteeContributionConsensusData = &types.ValidatorConsensusData{
	Duty:    TestingSyncCommitteeContributionDuty,
	DataSSZ: TestingContributionsDataBytes,
	Version: spec.DataVersionPhase0,
}
var TestSyncCommitteeContributionConsensusDataByts, _ = TestSyncCommitteeContributionConsensusData.Encode()
var TestSyncCommitteeContributionConsensusDataRoot = func() [32]byte {
	return sha256.Sum256(TestSyncCommitteeContributionConsensusDataByts)
}()
