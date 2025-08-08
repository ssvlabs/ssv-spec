package aggregatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SyncCommitteeContributionNoJustifications tests an invalid consensus data with no sync committee contribution pre-consensus justifications
func SyncCommitteeContributionNoJustifications() *AggregatorConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return NewAggregatorConsensusDataTest(
		"sync committee contribution with no pre-consensus justification",
		testdoc.AggregatorConsensusDataTestSyncCommitteeContributionNoJustificationsDoc,
		*testingutils.TestAggregatorCommitteeConsensusDataForDuty(testingutils.TestingSyncCommitteeContributionDuty, spec.DataVersionPhase0),
		"",
	)
}
