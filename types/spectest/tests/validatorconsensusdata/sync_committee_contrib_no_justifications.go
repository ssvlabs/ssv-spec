package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SyncCommitteeContributionNoJustifications tests an invalid consensus data with no sync committee contribution pre-consensus justifications
func SyncCommitteeContributionNoJustifications() *ValidatorConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	cd := types.ValidatorConsensusData{
		Duty:    testingutils.TestingSyncCommitteeContributionDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: testingutils.TestingContributionsDataBytes,
	}

	return &ValidatorConsensusDataTest{
		Name:          "sync committee contribution with no pre-consensus justification",
		ConsensusData: cd,
	}
}
