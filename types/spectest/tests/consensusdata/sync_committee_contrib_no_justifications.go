package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SyncCommitteeContributionNoJustifications tests an invalid consensus data with no sync committee contribution pre-consensus justifications
func SyncCommitteeContributionNoJustifications() *ConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	cd := types.ConsensusData{
		Duty:    testingutils.TestingSyncCommitteeContributionDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: testingutils.TestingContributionsDataBytes,
	}

	return &ConsensusDataTest{
		Name:          "sync committee contribution with no pre-consensus justification",
		ConsensusData: cd,
	}
}
