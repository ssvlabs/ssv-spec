package aggregatorcommitteeconsensusdata

import (
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SyncCommitteeContributionNoJustifications tests an invalid consensus data with no sync committee contribution pre-consensus justifications
func SyncCommitteeContributionNoJustifications() *AggregatorCommitteeConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	return NewAggregatorCommitteeConsensusDataTest(
		"sync committee contribution with no pre-consensus justification",
		testdoc.AggregatorCommitteeConsensusDataTestSyncCommitteeContributionNoJustificationsDoc,
		*testingutils.TestSyncCommitteeContributionConsensusData,
		0,
	)
}
