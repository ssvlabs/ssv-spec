package valcheckcommittee

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongAttestationWithConsensusData tests submitting an attestation ConsensusData which returns an error since it's not a BeaconVote object
func WrongAttestationWithConsensusData() tests.SpecTest {
	return &valcheck.SpecTest{
		Name:          "committee value check wrong attestation consensus data",
		Network:       types.BeaconTestNetwork,
		Role:          types.RoleCommittee,
		Input:         testingutils.TestAttesterConsensusDataByts,
		ExpectedError: "failed decoding beacon vote: incorrect size",
	}
}

// WrongSyncCommitteeWithConsensusData tests submitting a sync committee ConsensusData which returns an error since it's not a BeaconVote object
func WrongSyncCommitteeWithConsensusData() tests.SpecTest {
	return &valcheck.SpecTest{
		Name:          "committee value check wrong sync committee consensus data",
		Network:       types.BeaconTestNetwork,
		Role:          types.RoleCommittee,
		Input:         testingutils.TestSyncCommitteeConsensusDataByts,
		ExpectedError: "failed decoding beacon vote: incorrect size",
	}
}
